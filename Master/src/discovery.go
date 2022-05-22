package main

import (
	Type "DiniSQL/Region"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//ServiceDiscovery 服务发现 用于同步etcd的数据到本地
type ServiceDiscovery struct {
	cli           *clientv3.Client  //etcd client
	serverList    map[string]int    //连接上的Region, key是IP:port, value是Region的访问次数
	tableList     map[string]string // table name map to region
	regionCnt     int
	sortedRegions []RegionStatus
	lock          sync.Mutex
}

var endpoints = []string{"localhost:2379"}
var regionSer *ServiceDiscovery = NewServiceDiscovery(endpoints) // 每个region的IP:PORT和它对应的访问次数
//var tableSer *ServiceDiscovery = NewServiceDiscovery(endpoints)  // 每个table的拥有对应table的region们地址

//NewServiceDiscovery  新建发现服务
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]int),
		tableList:  make(map[string]string),
		regionCnt:  0,
	}
}

//WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService(prefix string) error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// 原有的k和v
	for _, ev := range resp.Kvs {
		s.UpdateList(string(ev.Key), string(ev.Value))
	}

	//监视前缀，修改变更的server
	go s.watcher(prefix)
	return nil
}

//watcher 监听前缀
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				// 监听到有新的region接入也是在这里
				s.UpdateList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

//UpdateList 用于通过etcd更新维护本地数据
func (s *ServiceDiscovery) UpdateList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if strings.HasPrefix(key, "/region/") {
		key = key[len("/region/"):]
		if value, ok := s.serverList[key]; ok {
			//int32 to string
			fmt.Println("region: " + key + " 已经存在, 原有访问次数为: " + strconv.Itoa(value) + ". 更新为: " + val)
		} else {
			// 有Region接入
			fmt.Println("region: " + key + " 首次连接, 访问次数为: " + val)
			regionSer.regionCnt++
			regionSer.sortedRegions = append(regionSer.sortedRegions, RegionStatus{
				regionIP:   strings.Split(key, ":")[0],
				regionPort: strings.Split(key, ":")[1],
				regionID:   0,
				rawStatus:  "",
				firstConn:  time.Now().Format("2006-01-02 15:04:05"),
				surviving:  true,
			})
		}
		s.serverList[key], _ = strconv.Atoi(val)
		log.Println("[ regionList PUT ]: ", "key :", key, "val:", val)
	} else if strings.HasPrefix(key, "/table/") {
		key = key[len("/table/"):]
		s.tableList[key] = val
		log.Println("[ tableList PUT ]: ", "key :", key, "val:", val)

	}
}

//DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if strings.HasPrefix(key, "/region/") {
		key = key[len("/region/"):]
		s.HandleRegionQuit(key)
		delete(s.serverList, key)
		log.Println("[ regionList DROP (region QUIT! )]: ", "key :", key)
	} else if strings.HasPrefix(key, "/table/") {
		key = key[len("/table/"):]
		delete(s.tableList, key)
		log.Println("[ tableList DROP ]: ", "key :", key)
	}
}

//HandleRegionQuit 处理某个region由于意外退出
// name 是region的IP:PORT
func (s *ServiceDiscovery) HandleRegionQuit(name string) {
	lostTables := s.tableList[name]
	for _, table := range strings.Split(lostTables, ";") {
		owndRegionList := s.tableList[table]
		hasTableRegionStr := strings.Replace(owndRegionList, name+";", "", -1)
		hasTableRegionList := strings.Split(hasTableRegionStr, ";")
		// 从拥有丢失表的region中挑一个发给另一个region
		srcRegion := hasTableRegionList[rand.Intn(len(hasTableRegionList))]
		targetRegion := domStr(regionSer.sortedRegions[0].regionIP, regionSer.sortedRegions[0].regionPort, true)
		p := Type.Packet{}
		p.Head = Type.PacketHead{P_Type: Type.UploadRegion, Op_Type: -1, Spare: ""}
		p.Payload = []byte(srcRegion + ";" + table + ";" + targetRegion) // TODO: implement real payload
		i, _ := strconv.Atoi(strings.Split(srcRegion, ":")[1])
		address := net.TCPAddr{
			IP:   net.ParseIP(strings.Split(srcRegion, ":")[0]),
			Port: i,
		}
		conn, err := net.DialTCP("tcp4", nil, &address)
		if err != nil {
			log.Fatal(err) // Println + os.Exit(1)
			return
		}
		var packetBuf = make([]byte, 0)
		packetBuf, err = p.MarshalMsg(packetBuf)
		_, err1 := conn.Write(packetBuf)
		if err1 != nil {
			log.Println(err)
			conn.Close()
			return
		}
		conn.Close()

	}
}

//GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)
	addrs = append(addrs, "region list->")
	for k, v := range s.serverList {
		addrs = append(addrs, fmt.Sprintf("%s:%d", k, v))
	}
	addrs = append(addrs, "table list->")
	for k, v := range s.tableList {
		addrs = append(addrs, fmt.Sprintf("%s:%s", k, v))
	}
	return addrs
}

//Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

func (s *ServiceDiscovery) HandleUpdate(regionSer *ServiceDiscovery, tableSer *ServiceDiscovery) {
	newRegionCnt := len(regionSer.serverList) // new server count
	if newRegionCnt < regionSer.regionCnt {
		log.Println("[ regionList DROP ]: ", "key :", regionSer.serverList)
		regionSer.regionCnt = newRegionCnt
	} else {
		regionSer.regionCnt = newRegionCnt
	}

}

func masterDiscovery() {
	//regionSer := NewServiceDiscovery(endpoints) // 每个region的IP:PORT和它对应的访问次数
	regionSer.WatchService("/region/")
	//tableSer := NewServiceDiscovery(endpoints) // 每个table的拥有对应table的region们地址
	regionSer.WatchService("/table/")
	defer regionSer.Close()
	//regionSer.WatchService("/web/")
	//regionSer.WatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(5 * time.Second):
			log.Println("Region count: ", regionSer.regionCnt)
			log.Println(regionSer.GetServices())
		}
	}
}

package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli        *clientv3.Client  //etcd client
	serverList map[string]string //连接上的Region, key是IP:port, value是Region的访问次数
	tableList  map[string]string
	lock       sync.Mutex
}

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
		serverList: make(map[string]string),
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
				s.UpdateList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

//UpdateList 新增服务地址
func (s *ServiceDiscovery) UpdateList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if strings.HasPrefix(key, "/region/") {
		key = key[len("/region/"):]
		s.serverList[key] = val
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
		delete(s.serverList, key)
		log.Println("[ regionList DROP ]: ", "key :", key)
	} else if strings.HasPrefix(key, "/table/") {
		key = key[len("/table/"):]
		delete(s.tableList, key)
		log.Println("[ tableList DROP ]: ", "key :", key)
	}
}

//GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *ServiceDiscovery) updRegStat() {
	sortedRegions = []RegionStatus{}
}

//Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

func masterDiscovery() {
	var endpoints = []string{"localhost:2379"}
	regionSer := NewServiceDiscovery(endpoints) // 每个region的IP:PORT和它对应的访问次数
	regionSer.WatchService("/region/")
	tableSer := NewServiceDiscovery(endpoints) // 每个table的拥有对应table的region们地址
	tableSer.WatchService("/table/")
	defer regionSer.Close()
	defer tableSer.Close()
	//regionSer.WatchService("/web/")
	//regionSer.WatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(5 * time.Second):
			// TODO:维护regionStatus
			// TODO:维护sortedRegions
			regionCnt = len(regionSer.serverList) // update server count

			log.Println(regionSer.GetServices())
		}
	}
}

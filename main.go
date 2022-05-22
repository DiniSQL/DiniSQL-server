package main

import (
	"DiniSQL/MiniSQL"
	"DiniSQL/Region"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
)

func runEtcd() {
	// cmd := exec.Command("etcd", "--name", "infra1", "--initial-advertise-peer-urls", "http://192.168.84.244:2380",
	// 	"--listen-peer-urls", "http://192.168.84.244:2380",
	// 	"--listen-client-urls", "http://192.168.84.244:2379,http://127.0.0.1:2379",
	// 	"--advertise-client-urls", "http://192.168.84.244:2379",
	// 	"--initial-cluster-token", "etcd-cluster-1",
	// 	"--initial-cluster", "infra0=http://192.168.84.48:2380,infra1=http://192.168.84.244:2380",
	// 	"--initial-cluster-state", "new")
	cmd := exec.Command("etcd")
	stout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	defer stout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	if opBytes, err := ioutil.ReadAll(stout); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(string(opBytes))
	}

}

// getLocalIpV4 获取 IPV4 IP，没有则返回空
func getLocalIpV4(interfaceName string) (addr string, err error) {
	inter, err := net.InterfaceByName(interfaceName)
	if err != nil {
		panic(err)
	}
	// 判断网卡是否开启，过滤本地环回接口
	if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
		// 获取网卡下所有的地址
		addrs, err := inter.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				//判断是否存在IPV4 IP 如果没有过滤
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("interface %s don't have an ipv4 address", interfaceName)
}

func main() {
	ip, err := getLocalIpV4("wifi0")
	if err != nil {
		log.Fatal(err)
	}

	// runEtcd()
	var endpoints = []string{"127.0.0.1:2379"}
	prefix := "/region"
	key := strings.Join([]string{prefix, ip}, "/")
	fmt.Println(key)
	ser, err := Region.NewServiceRegister(endpoints, key, "0", 5)
	if err != nil {
		log.Fatalln(err)
	}
	//监听续租相应chan
	go ser.ListenLeaseRespChan()

	MiniSQL.InitDB()
	Region.InitRegionServer()
}

// func main() {
// 	//errChan 用于接收shell返回的err
// 	errChan := make(chan error)
// 	go MiniSQL.RunShell(errChan) //开启shell协程
// 	err := <-errChan
// 	fmt.Println("bye")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

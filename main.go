package main

import (
	"DiniSQL/MiniSQL"
	"DiniSQL/Region"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func runEtcd() {
	cmd := exec.Command("etcd", "--name", "infra1", "--initial-advertise-peer-urls", "http://192.168.84.244:2380",
		"--listen-peer-urls http://192.168.84.244:2380",
		"--listen-client-urls http://192.168.84.244:2379,http://127.0.0.1:2379",
		"--advertise-client-urls http://192.168.84.244:2379",
		"--initial-cluster-token etcd-cluster-1",
		"--initial-cluster infra0=http://192.168.84.48:2380,infra1=http://192.168.84.244:2380",
		"--initial-cluster-state new")
	// cmd := exec.Command("etcd", "--version")
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

// func main() {
// 	// runEtcd()
// 	var endpoints = []string{"127.0.0.1:2379"}
// 	ser, err := Region.NewServiceRegister(endpoints, "/web/node1", "localhost:8000", 5)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//监听续租相应chan
// 	go ser.ListenLeaseRespChan()
// 	select {
// 	// case <-time.After(20 * time.Second):
// 	// 	ser.Close()
// 	}
// }

func main() {
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

// func main() {
// 	var endpoints = []string{"127.0.0.1:2379"}
// 	ser, err := Region.NewServiceRegister(endpoints, "/web/node1", "localhost:8000", 5)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//监听续租相应chan
// 	go ser.ListenLeaseRespChan()
// 	select {
// 	// case <-time.After(20 * time.Second):
// 	// 	ser.Close()
// 	}
// }

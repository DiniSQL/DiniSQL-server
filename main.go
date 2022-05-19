package main

import (
<<<<<<< HEAD
	"DiniSQL/MiniSQL/src/API"
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/CatalogManager"
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"DiniSQL/MiniSQL/src/RecordManager"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterh/liner"
)

const historyCommmandFile = "~/.DiniSQL_history"

//FlushALl 结束时做的所有工作
func FlushALl() {
	BufferManager.BlockFlushAll()                                             //缓存block
	RecordManager.FlushFreeList()                                             //free list写回
	CatalogManager.FlushDatabaseMeta(CatalogManager.UsingDatabase.DatabaseId) //刷新记录长度和余量
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		parts := strings.SplitN(path, "/", 2)
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, parts[1]), nil
=======
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
>>>>>>> d83281e932449187e7dedc8b7d5e4db7c4847aaa
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
<<<<<<< HEAD

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
=======
>>>>>>> d83281e932449187e7dedc8b7d5e4db7c4847aaa

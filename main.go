package main

import (
	"DiniSQL/Region"
	"fmt"
)

func main() {
	var ip string
	fmt.Scanln(&ip)
	Region.InitRegionServer(ip)
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

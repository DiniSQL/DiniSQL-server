package main

import (
	"fmt"
	// "github.com/peterh/liner"
)
func start(){
	fmt.Println("MiniSQL客户端启动！")
	fmt.Print("请输入:")
	var input string
	fmt.Scan(&input)
	for input!=`exit`{
		fmt.Println("输入："+input)
		fmt.Print("请输入:")
		fmt.Scan(&input)
	}
	fmt.Println("输入:"+input)
	// var input string

}
func main() {
    // start()
	// ll := liner.NewLiner()
	// defer ll.Close()
	// ll.SetCtrlCAborts(true)
	newCache()
	setCache("table1","region1")
	setCache("table2","region2")
	fmt.Println("getCache1"+getCache("table1"))
	fmt.Println(cache)
	fmt.Println("getCache3"+getCache("table3"))
	flushCache()
	fmt.Println(cache)

}
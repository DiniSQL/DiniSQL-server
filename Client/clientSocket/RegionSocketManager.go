package main

import (
    "fmt"
    "log"
    "net"

)

var regionIP string = "localhost"
var regionPort int = 12345

// initial connection to region server and send message
func ConnectToRegionIP(regionIP string, sql string)bool{
    address := net.TCPAddr{
        IP:   net.ParseIP(regionIP), // 把字符串IP地址转换为net.IP类型
        Port: regionPort,
    }
	conn, err := net.DialTCP("tcp4", nil, &address)
    if err != nil {
        log.Fatal(err) // Println + os.Exit(1)
		return false
    }
	_, err1 := conn.Write([]byte(sql))
	if err1 != nil {
		log.Println(err)
		conn.Close()
		return false
	}
	fmt.Printf("send %s to %s\n", sql, conn.RemoteAddr())
	return true
}
// initial connection to region server and send message
func ConnectToRegionPort(regionPort int, sql string)bool{
    address := net.TCPAddr{
        IP:   net.ParseIP(regionIP), // 把字符串IP地址转换为net.IP类型
        Port: regionPort,
    }
	conn, err := net.DialTCP("tcp4", nil, &address)
    if err != nil {
        log.Fatal(err) // Println + os.Exit(1)
		return false
    }
	_, err1 := conn.Write([]byte(sql))
	if err1 != nil {
		log.Println(err)
		conn.Close()
		return false
	}
	fmt.Printf("send %s to %s\n", sql, conn.RemoteAddr())
	conn.Close()
	return true
}

// listen 
// input : IP and Port of client
func KeepListening(ClientIP string, ClientPort int){
    address := net.TCPAddr{
        IP:   net.ParseIP(ClientIP), // 把字符串IP地址转换为net.IP类型
        Port: ClientPort,
    }
    listener, err := net.ListenTCP("tcp4", &address) // 创建TCP4服务器端监听器
    if err != nil {
        log.Fatal(err) // Println + os.Exit(1)
    }
    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            log.Fatal(err) // 错误直接退出
        }
        fmt.Println("remote address:", conn.RemoteAddr())
        res := make([]byte,1024)
        _, err1 := conn.Read(res)
		if err1 != nil {
            log.Println(err)
            conn.Close()
        }
        fmt.Printf("received %s from %s", res, conn.RemoteAddr())
    }
}

func main(){
	var sql string
	go KeepListening("127.0.0.1",8004)
	for true  {
        fmt.Scanln(&sql)
		ConnectToRegionPort(8000,sql)
    }
}
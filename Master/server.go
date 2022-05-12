package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func echo(conn *net.TCPConn) {
    tick := time.Tick(5 * time.Second) // 五秒的心跳间隔
    for now := range tick {
        n, err := conn.Write([]byte(now.String()))
        fmt.Println(now)
        if err != nil {
            log.Println(err)
            conn.Close()
            return
        }
        fmt.Printf("send %d bytes to %s\n", n, conn.RemoteAddr())
    }
}
func main() {
    address := net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"), // 把字符串IP地址转换为net.IP类型
        Port: 8000,
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
        s := make([]byte, 1024)
        // var s [1024]byte
        n,_ := conn.Read(s)
        fmt.Printf("----%d----%s----",n,s)
        fmt.Println("remote address:", conn.RemoteAddr())
        go echo(conn)
    }
}
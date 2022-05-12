package main

import (
    "log"
    "net"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        log.Fatalf("Usage: %s host:port", os.Args[0])
    }
    service := os.Args[1]
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    if err != nil {
        log.Fatal(err)
    }
    conn, err := net.DialTCP("tcp4", nil, tcpAddr)
    if err != nil {
        log.Fatal(err)
    }
    _,err1 := conn.Write([]byte("123456"))
    if err1 != nil {
        log.Fatal(err)
    }
    // log.Fatal(n)
}
package Region

import (
	"fmt"
	"net"
)

var regionServer *RegionServer

type RegionServer struct {
	regions    []region
	serverID   int
	visitCount int
}

func initRegionServer() {
	regionServer = new(RegionServer)
	regionServer.visitCount = 0
	listenFromClient(regionServer)
}

func listenFromClient(server *RegionServer) {
	listen, err := net.Listen("tcp", "")
	if err != nil {
		fmt.Println("Failed to listen to client!")
		return
	}
	for {
		conn, err := listen.Accept()
		fmt.Println("Accepted an connection")
		if err != nil {
			fmt.Println("Accept failed!")
			continue
		}
		go server.serveClient(conn)
	}
}

func (server *RegionServer) serveClient(conn net.Conn) {
	defer conn.Close()
	server.visitCount++
	var buf = make([]byte, 1024)
	p := Packet{}
	for {
		conn.Read(buf)
		p.UnmarshalMsg(buf)
	}
}

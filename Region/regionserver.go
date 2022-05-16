package Region

import (
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"fmt"
	"net"
	"time"
	"unsafe"
)

var regionServer *RegionServer

type RegionServer struct {
	regions    []region
	serverID   int
	visitCount int
}

func InitRegionServer() {
	regionServer = new(RegionServer)
	regionServer.visitCount = 0

	listenFromClient(regionServer)
}

func (server *RegionServer) heartBeat(conn net.Conn) {
	for {
		time.Sleep(5 * time.Second)
		str := fmt.Sprintf("ServerID:%d", server.serverID)
		//h := HeartBeat2etcd{serverID: server.serverID, regions: regions}
		p := Packet{Head: PacketHead{P_Type: KeepAlive, Op_Type: -1},
			Payload: []byte(str)}
		var replyBuf = make([]byte, p.Msgsize())
		p.MarshalMsg(replyBuf)
		conn.Write(replyBuf)
	}
}

func listenFromClient(server *RegionServer) {
	listen, err := net.Listen("tcp", "172.20.10.3")
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
		go server.serve(conn)
	}
}

func byteSliceToString(bytes []byte) string {

	return *(*string)(unsafe.Pointer(&bytes))

}

func (server *RegionServer) serve(conn net.Conn) {
	defer conn.Close()
	server.visitCount++
	p := Packet{}
	for {
		var buf = make([]byte, 1024)
		conn.Read(buf)
		var content = byteSliceToString(buf)
		if content == "end" {
			break
		} else {
			p.UnmarshalMsg(buf)
			if p.Head.P_Type == SQLOperation {
				var statement = byteSliceToString(p.Payload)
				switch p.Head.Op_Type {
				case types.Select:
					fmt.Println(statement)
				case types.CreateDatabase:
					fmt.Println(statement)
				case types.DropDatabase:
					fmt.Println(statement)
				case types.CreateIndex:
					fmt.Println(statement)
				case types.DropIndex:
					fmt.Println(statement)
				case types.Insert:
					fmt.Println(statement)
				default:
					fmt.Println(statement)
				}
			}
			replyPacket := Packet{Head: PacketHead{P_Type: Result, Op_Type: -1},
				Payload: []byte{1, 1, 1}}
			var replyBuf = make([]byte, replyPacket.Msgsize())
			replyPacket.MarshalMsg(replyBuf)
			conn.Write(replyBuf)
		}
	}
}

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
		regions := make([]simpleRegion, len(server.regions))
		for i := 1; i < len(server.regions); i++ {
			regions[i] = server.regions[i].simplify()
		}
		str := fmt.Sprintf("ServerID:%d", server.serverID)
		//h := HeartBeat2etcd{serverID: server.serverID, regions: regions}
		p := Packet{head: PacketHead{packetType: KeepAlive, detailedType: -1},
			payload: Payload{content: []byte(str)}}
		var replyBuf = make([]byte, p.Msgsize())
		p.MarshalMsg(replyBuf)
		conn.Write(replyBuf)
	}
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
			if p.head.packetType == SQLOperation {
				var statement = byteSliceToString(p.payload.content)
				switch p.head.detailedType {
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
			replyPacket := Packet{head: PacketHead{packetType: Result, detailedType: -1}, payload: Payload{content: []byte{1, 1, 1}}}
			var replyBuf = make([]byte, replyPacket.Msgsize())
			replyPacket.MarshalMsg(replyBuf)
			conn.Write(replyBuf)
		}
	}
}

package Region

import (
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"DiniSQL/MiniSQL/src/Utils/Error"
	"DiniSQL/MinisQL/src/API"
	"fmt"
	"net"
	"strings"
	"time"
	"unsafe"
)

var regionServer *RegionServer
var StatementChannel chan types.DStatements
var FinishChannel chan string
var FlushChannel chan struct{}

type RegionServer struct {
	regions    []region
	serverID   int
	visitCount int
}

func InitRegionServer() {
	regionServer = new(RegionServer)
	regionServer.visitCount = 0
	StatementChannel = make(chan types.DStatements, 500)   //用于传输操作指令通道
	FinishChannel = make(chan Error.Error, 500)            //用于api执行完成反馈通道
	FlushChannel = make(chan struct{})                     //用于每条指令结束后协程flush
	go API.HandleOneParse(StatementChannel, FinishChannel) //begin the runtime for exec
	go BufferManager.BeginBlockFlush(FlushChannel)
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
		replyBuf, err := p.MarshalMsg(replyBuf)
		if err != nil {
			fmt.Println(err)
		}
		conn.Write(replyBuf)
	}
}

func listenFromClient(server *RegionServer) {
	listen, err := net.Listen("tcp", "127.0.0.1:3037")
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
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error: conn.Read()")
			break
		}
		var content = byteSliceToString(buf)
		if content == "end" {
			break
		} else {
			p.UnmarshalMsg(buf)
			var res string
			if p.Head.P_Type == SQLOperation {
				var statement = byteSliceToString(p.Payload)
				err = parser.Parse(strings.NewReader(string(statement)), StatementChannel)
				res = <-FinishChannel
			}
			replyPacket := Packet{Head: PacketHead{P_Type: Result, Op_Type: -1},
				Payload: []byte(res)}
			var replyBuf = make([]byte, 0)
			replyBuf, err = replyPacket.MarshalMsg(replyBuf)
			if err != nil {
				fmt.Println(err)
			}
			_, err := conn.Write(replyBuf)
			if err != nil {
				fmt.Println("Error: conn.Write()")
				break
			}
		}
	}
}

func Select(statement string) {

}

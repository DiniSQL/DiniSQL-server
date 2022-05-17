package Region

import (
	"DiniSQL/MiniSQL/src/API"
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"fmt"
	"log"
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
	FinishChannel = make(chan string, 500)                 //用于api执行完成反馈通道
	FlushChannel = make(chan struct{})                     //用于每条指令结束后协程flush
	go API.HandleOneParse(StatementChannel, FinishChannel) //begin the runtime for exec
	go BufferManager.BeginBlockFlush(FlushChannel)
	err := parser.Parse(strings.NewReader(string("create database minisql;")), StatementChannel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(<-FinishChannel)
	err = parser.Parse(strings.NewReader(string("use database minisql;")), StatementChannel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(<-FinishChannel)
	// err = parser.Parse(strings.NewReader(string("insert into student values(1080100001,'name1',99);")), StatementChannel)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(<-FinishChannel)
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
	listen, err := net.Listen("tcp", ":3037")
	fmt.Println("Server started")
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
	var p Packet

	var res = make([]byte, 1024)
	cnt, err := conn.Read(res)
	fmt.Println(cnt)
	if err != nil {
		log.Println(err)
		conn.Close()
	}
	_, err = p.UnmarshalMsg(res)
	if err != nil {
		fmt.Println("UnmarshalMsg failed")
		fmt.Println(err)
	}
	fmt.Println(p.Head.P_Type)
	var opRes string
	if p.Head.P_Type == SQLOperation {
		var statement = byteSliceToString(p.Payload)
		fmt.Println(statement)
		err = parser.Parse(strings.NewReader(string(statement)), StatementChannel)
		if err != nil {
			fmt.Println(err)
		}
		opRes = <-FinishChannel
		fmt.Println(opRes)
	} else if p.Head.P_Type == End {
	}
	replyPacket := Packet{Head: PacketHead{P_Type: Result, Op_Type: -1},
		Payload: []byte(opRes)}
	var replyBuf = make([]byte, 0)
	replyBuf, err = replyPacket.MarshalMsg(replyBuf)
	if err != nil {
		fmt.Println("MarshalMsg failed")
	}
	_, err = conn.Write(replyBuf)
	if err != nil {
		fmt.Println("Error: conn.Write()")
	}
}

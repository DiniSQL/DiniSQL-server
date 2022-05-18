package main

import (
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	Type "DiniSQL/Region"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

type RegionStatus struct {
	regionIP   string
	regionPort string
	regionID   int
	rawStatus  string
	lastConn   int64
	surviving  bool
	loadLevel  int
}

const (
	MASTER_PORT = ":9000"
	ClientIP    = "127.0.0.1"
	ClientPort  = 6100
)

var tab2reg map[string]*RegionStatus
var id2reg map[int]*RegionStatus

func initMaster() {
	// Initialize the master server
	tab2reg = make(map[string]*RegionStatus) // table name to region status
	tab2reg["foo"] = new(RegionStatus)
	id2reg = make(map[int]*RegionStatus) // region id to region status
	id2reg[0] = new(RegionStatus)
	id2reg[0].regionID = 1

}

func ConnectToRegion(regionIP string, regionPort int, packet Type.Packet) (recPacket Type.Packet) {

	address := net.TCPAddr{
		IP:   net.ParseIP(regionIP),
		Port: regionPort,
	}
	conn, err := net.DialTCP("tcp4", nil, &address)
	// defer conn.Close()
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
		return
	}
	var packetBuf = make([]byte, 0)
	fmt.Printf("packet.Head.P_Type:%d\n", packet.Head.P_Type)
	fmt.Printf("packet.Head.Op_Type:%d\n", packet.Head.Op_Type)
	fmt.Printf("packet.Payload:%s\n", packet.Payload)
	packetBuf, err = packet.MarshalMsg(packetBuf)
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
		return
	}
	_, err1 := conn.Write(packetBuf)
	fmt.Println(packetBuf)
	if err1 != nil {
		log.Println(err)
		conn.Close()
		return
	}
	conn.Close()
	fmt.Printf("send %d to %s\n", packet.Head.P_Type, conn.RemoteAddr())
	// listen after send
	recPacket = KeepListening(ClientIP, ClientPort)
	return

}

func ConnectToClient(regionIP string, regionPort int, packet Type.Packet) (recPacket Type.Packet) {

	address := net.TCPAddr{
		IP:   net.ParseIP(regionIP),
		Port: regionPort,
	}
	conn, err := net.DialTCP("tcp4", nil, &address)
	//conn, err := net.Dial("tcp", "172.20.10.3:8005")
	// defer conn.Close()
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
		return
	}
	var packetBuf = make([]byte, 0)
	//fmt.Printf("packet.Head.P_Type:%d\n", packet.Head.P_Type)
	//fmt.Printf("packet.Head.Op_Type:%d\n", packet.Head.Op_Type)
	//fmt.Printf("packet.Payload:%s\n", packet.Payload)
	packetBuf, err = packet.MarshalMsg(packetBuf)
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
		return
	}
	_, err1 := conn.Write(packetBuf)
	fmt.Println(packetBuf)
	if err1 != nil {
		log.Println(err)
		conn.Close()
		return
	}
	conn.Close()
	fmt.Printf("send %d to %s\n", packet.Head.P_Type, conn.RemoteAddr())
	// listen after send
	//recPacket = KeepListening(ClientIP, ClientPort)
	return

}

// listen
// input : IP and Port of client
func KeepListening(ClientIP string, ClientPort int) (receivedPacket Type.Packet) {
	fmt.Println("listening...")
	listener, err := net.Listen("tcp", MASTER_PORT)
	defer listener.Close()
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("remote address:", conn.RemoteAddr())
	res, err1 := ioutil.ReadAll(conn)
	if err1 != nil {
		log.Println(err)
		conn.Close()
	}
	res, err1 = receivedPacket.UnmarshalMsg(res)
	// ch <- p
	if err1 != nil {
		log.Println(err)
		conn.Close()
	}

	fmt.Printf("p.Head.P_Type:%d\n", receivedPacket.Head.P_Type)
	fmt.Printf("p.Head.Op_Type:%d\n", receivedPacket.Head.Op_Type)
	fmt.Printf("p.Payload:%s\n", receivedPacket.Payload)

	//p := Type.Packet{Head: Type.PacketHead{P_Type: Type.KeepAlive, Op_Type: Type.CreateIndex},
	//	Payload: []byte("foo")}
	//fmt.Printf("remote addr: %s\n", conn.RemoteAddr().String())
	//remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")
	//ConnectToClient(remoteAddr[0], 8005, p)

	return

}

func main() {
	initMaster()
	//reader, err := os.Open("1.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	StatementChannel := make(chan types.DStatements, 500)
	FinishChannel := make(chan string, 500)
	TablesChannel := make(chan []string, 500)
	//FlushChannel := make(chan struct{})
	go Parse2Statement(StatementChannel, FinishChannel, TablesChannel)
	//fmt.Println("Initialized Master server")
	var sql_strings = []string{
		"create table tab1(a int);",
		"select a from tab2;",
	}

	err := parser.Parse(strings.NewReader(sql_strings[0]), StatementChannel) //开始解析
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(<-TablesChannel)
	close(StatementChannel) //关闭StatementChannel，进而关闭FinishChannel
	for _ = range FinishChannel {

	}
	for {
		client := KeepListening(ClientIP, 9000)
		fmt.Println(client)

		//ConnectToRegion("127.0.0.1", 8006, p)
	}

}

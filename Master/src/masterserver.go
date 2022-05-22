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
	firstConn  string
	surviving  bool
}

const (
	MASTER_PORT = ":9000"
	ClientIP    = "127.0.0.1"
	ClientPort  = 6100
)

var tab2reg map[string][]RegionStatus

var StatementChannel chan types.DStatements
var FinishChannel chan string
var TablesChannel chan []string
var OutputStatementChannel chan types.DStatements

//var regionStatus map[string]*RegionStatus // region name to region status

//var StatementChannel = make(chan types.DStatements, 500)

func initMaster() {
	// Initialize the master server
	tab2reg = make(map[string][]RegionStatus) // table name to region status
	tab2reg["foo"] = append(tab2reg["foo"], RegionStatus{})
	regionSer.regionCnt = 0
	//id2reg = make(map[int]*RegionStatus) // region id to region status
	//id2reg[0] = new(RegionStatus)
	//id2reg[0].regionID = 1

}

func domStr(IP string, port string, end bool) string {
	if end {
		return IP + ":" + port + ";"
	} else {
		return IP + ":" + port
	}
}

func regionList(status []RegionStatus) string {
	var ret string
	for _, region := range status {
		ret += domStr(region.regionIP, region.regionPort, true)
	}
	return ret
}

func findRelaxRegion(count int) string {
	var ret string
	selectedRegion := 0
	targetCount := 0 // 当前遍历到的访问次数，从最低的0开始
	// 假设我们所需的服务器数量count永远能小于总连接的region
	for {
		if selectedRegion >= count {
			break
		}
		for region, currentCount := range regionSer.serverList {
			if targetCount == currentCount {
				ret += region + ";"
				selectedRegion++
			}
		}
		// 每个循环中，所寻找的指定访问次数++
		targetCount++

	}

	return ret

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
	recPacket = HandleClient(ClientIP, ClientPort)
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
	//recPacket = HandleClient(ClientIP, ClientPort)
	return

}

func CreatePacket(statement types.DStatements, tables []string) Type.Packet {
	packet := Type.Packet{}
	packet.Head = Type.PacketHead{P_Type: Type.Answer, Op_Type: statement.GetOperationType(), Spare: ""}
	var pay string
	switch statement.GetOperationType() {
	case types.CreateDatabase:
	case types.UseDatabase:
	case types.CreateTable:
		selCnt := min(2, regionSer.regionCnt)   // select specified number of regions
		relaxRegions := findRelaxRegion(selCnt) //a.a.a:123;b.b.b:456;c.c.c:789;
		pay = relaxRegions
		//pay = regionList(relaxRegions)
		//tab2reg[tables[0]] = relaxRegions
		regionSer.tableList[tables[0]] = relaxRegions

	case types.CreateIndex:
		pay = regionSer.tableList[tables[0]]

	case types.DropTable:
		pay = regionSer.tableList[tables[0]]

	case types.DropIndex:
		pay = regionSer.tableList[tables[0]]

	case types.Insert:
		pay = regionSer.tableList[tables[0]]

	case types.Update:
		pay = regionSer.tableList[tables[0]]

	case types.Delete:
		pay = regionSer.tableList[tables[0]]

	case types.Select:
		pay = regionSer.tableList[tables[0]]

	case types.ExecFile:
	case types.DropDatabase:

	}
	packet.Payload = []byte(pay)
	return packet
}

func HandleClient(ClientIP string, ClientPort int) (receivedPacket Type.Packet) {
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
	fmt.Println("remote connected! address:", conn.RemoteAddr())
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
	err = parser.Parse(strings.NewReader(string(receivedPacket.Payload)), StatementChannel) //收到客户端传过来的语句
	if err != nil {
		log.Fatal(err)
	}
	tabs := <-TablesChannel
	statement := <-StatementChannel
	fmt.Println("tabs:", tabs)

	p := CreatePacket(statement, tabs)
	packetBuf := make([]byte, 500)
	packetBuf, err = p.MarshalMsg(packetBuf)
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
		return
	}
	_, err1 = conn.Write(packetBuf)
	//fmt.Printf("remote addr: %s\n", conn.RemoteAddr().String())
	//remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")
	//ConnectToClient(remoteAddr[0], 8005, p)

	return

}

// 监听client的连接，当某个client连接上之后，处理其连接，等待请求的表的信息，获取表名之后将对应的Region地址返回给client
// 如果是创建新表，需要选择一个region，通知region新建表项
// 对于不同的
func main() {
	initMaster()

	StatementChannel = make(chan types.DStatements, 500)
	FinishChannel = make(chan string, 500)
	TablesChannel = make(chan []string, 500)
	OutputStatementChannel = make(chan types.DStatements, 500)
	//FlushChannel := make(chan struct{})
	go Parse2Statement(StatementChannel, FinishChannel, TablesChannel, OutputStatementChannel)
	//fmt.Println("Initialized Master server")
	var sql_strings = []string{
		"create table tab1(a int);",
		"select a from tab2;",
	}

	err := parser.Parse(strings.NewReader(sql_strings[0]), StatementChannel) //开始解析
	if err != nil {
		log.Fatal(err)
	}
	go masterDiscovery()
	//CreatePacket(<-OutputStatementChannel, <-TablesChannel)
	//fmt.Println(<-OutputStatementChannel, <-TablesChannel)
	close(StatementChannel) //关闭StatementChannel，进而关闭FinishChannel
	for _ = range FinishChannel {

	}
	for {
		client := HandleClient(ClientIP, 9000)
		fmt.Println(client)
	}

}

package main

import (
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	Type "DiniSQL/Region"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/tinylib/msgp/msgp"
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
	ClientPort  = 9000
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
			fmt.Println("region:", region, "currentCount:", currentCount)
			if targetCount == currentCount {
				ret += region + ";"
				selectedRegion++
			}
			if selectedRegion >= count {
				break
			}
		}
		// 每个循环中，所寻找的指定访问次数++
		targetCount++

	}
	return ret
}

func findTargetRegion(IPList string) (IP string, Port string) {
	Regions := strings.Split(IPList, ";")
	lowestLoad := -1
	for _, region := range Regions {
		if regionSer.serverList[region] < lowestLoad || lowestLoad == -1 {
			lowestLoad = regionSer.serverList[region]
			IP, Port = strings.Split(region, ":")[0], strings.Split(region, ":")[1]
		}
	}
	return
}

func ConnectToRegion(regionIP string, regionPort string, packet Type.Packet) (recPacket Type.Packet) {
	// 找到负载最低的region并连接它
	//address := net.TCPAddr{
	//	IP:   net.ParseIP(regionIP),
	//	Port: regionPort,
	//}
	//conn, err := net.DialTCP("tcp4", nil, &address)
	conn, err := net.Dial("tcp", domStr(regionIP, regionPort, false))
	// defer func(conn net.Conn) {
	// 	err := conn.Close()
	// 	if err != nil {
	// 		log.Println("Error when closing connection:", err)
	// 	}
	// }(conn)
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
		return
	}
	//rd := msgp.NewReader(conn)
	packet.IPResult = make([]byte, 0)

	//wt := msgp.NewWriter(conn)
	if err != nil {
		log.Println(err)
	}

	var packetBuf = make([]byte, 0)
	packetBuf, err = packet.MarshalMsg(packetBuf)
	_, err = conn.Write(packetBuf)
	//err = packet.EncodeMsg(wt)
	if err != nil {
		log.Println(err)
	}

	//fmt.Printf("packet.Head.P_Type:%d\n", packet.Head.P_Type)
	//fmt.Printf("packet.Head.Op_Type:%d\n", packet.Head.Op_Type)
	//fmt.Printf("packet.Payload:%s\n", packet.Payload)

	fmt.Printf("send %d to %s\n", packet.Head.P_Type, conn.RemoteAddr())

	// 获取region的返回包
	//err = recPacket.DecodeMsg(rd)
	packetBuf, err = ioutil.ReadAll(conn)
	//for {
	//	if len(packetBuf) != 0 {
	//		break
	//	}
	//
	//}
	if err != nil {
		log.Fatal(err)
	}
	_, err = recPacket.UnmarshalMsg(packetBuf)
	if err != nil {
		log.Fatal(err)
	}

	return

}

func SendToRegions(IPList string, packet Type.Packet) (recPacket Type.Packet) {
	regions := strings.Split(IPList, ";")
	finalStatus := true
	var finalResult []byte
	for _, region := range regions[:len(regions)-1] {
		IP, Port := strings.Split(region, ":")[0], strings.Split(region, ":")[1]
		//p_int, err := strconv.Atoi(Port)
		//if err != nil {
		//	return
		//}
		recvPkt := ConnectToRegion(IP, Port, packet)
		finalStatus = recvPkt.Signal
		finalResult = recvPkt.Payload
		if !finalStatus {
			break
		}
	}
	recPacket.Head = Type.PacketHead{P_Type: Type.Answer, Op_Type: 0, Spare: ""}
	recPacket.Payload = finalResult
	recPacket.Signal = finalStatus

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

func CreatePacket(statement types.DStatements, tables []string, SQLContent string) Type.Packet {
	//regionSer.lock.Lock()
	//defer regionSer.lock.Unlock()
	packetToClient := Type.Packet{}
	packetToClient.Head = Type.PacketHead{P_Type: Type.Answer, Op_Type: statement.GetOperationType(), Spare: ""}

	// 发给从节点对时候需要请求的是SQL的语句
	packetToRegion := Type.Packet{}
	packetToRegion.Head = Type.PacketHead{P_Type: Type.SQLOperation, Op_Type: statement.GetOperationType(), Spare: ""}
	packetToRegion.Payload = []byte(SQLContent)
	var IPResult string
	recvResult := Type.Packet{}

	IPResult = regionSer.tableList[tables[0]]
	switch statement.GetOperationType() {
	case types.CreateDatabase:

	case types.UseDatabase:
	case types.CreateTable:
		selCnt := min(2, regionSer.regionCnt) // select specified number of regions
		selCnt = 2                            // debug

		if selCnt == 0 {
			log.Fatal("no region available")
		}
		relaxRegions := findRelaxRegion(selCnt) //a.a.a:123;b.b.b:456;c.c.c:789;
		fmt.Println("selected regions: ", relaxRegions)
		IPResult = relaxRegions
		//IPResult = regionList(relaxRegions)
		//tab2reg[tables[0]] = relaxRegions
		regionSer.tableList[tables[0]] = relaxRegions
		recvResult = SendToRegions(IPResult, packetToRegion)
		regionSer.cli.Put(context.Background(), "/table/"+tables[0], IPResult)

	case types.CreateIndex:
		recvResult = SendToRegions(IPResult, packetToRegion)

	case types.DropTable:
		recvResult = SendToRegions(IPResult, packetToRegion)
		regionSer.cli.Delete(context.Background(), "/table/"+tables[0])
		delete(regionSer.tableList, tables[0])
		//oldRegions := regionSer.tableList[tables[0]]
		//oldRegions = strings.Replace(oldRegions, region+";", "", 1)
		//regionSer.tableList[tables[0]] = oldRegions
		//regionSer.cli.Put(context.Background(), "/table/"+tableName, regions)

	case types.DropIndex:
		recvResult = SendToRegions(IPResult, packetToRegion)

	case types.Insert:
		recvResult = SendToRegions(IPResult, packetToRegion)

	case types.Update:
		recvResult = SendToRegions(IPResult, packetToRegion)

	case types.Delete:
		recvResult = SendToRegions(IPResult, packetToRegion)

	case types.Select:
		// select语句不需要发送到region
		//recvResult = SendToRegions(IPResult, packetToRegion)

	case types.ExecFile:
	case types.DropDatabase:

	}
	packetToClient.IPResult = []byte(IPResult)
	packetToClient.Payload = []byte(recvResult.Payload)
	packetToClient.Signal = recvResult.Signal
	if statement.GetOperationType() == types.Select {
		packetToClient.Signal = true
	}
	//packetToClient.Payload = []byte("12345")
	return packetToClient
}

func HandleClient(conn net.Conn) (receivedPacket Type.Packet) {
	//if err != nil {
	//	log.Fatal(err)
	//}
	defer conn.Close()
	fmt.Println("remote connected! address:", conn.RemoteAddr())
	rd := msgp.NewReader(conn)
	//wt := msgp.NewWriter(conn)

	err := receivedPacket.DecodeMsg(rd)
	if err != nil {
		log.Fatal(err)
	}
	//res, err1 := ioutil.ReadAll(conn)
	//if err1 != nil {
	//	log.Println(err)
	//	conn.Close()
	//}
	//res, err1 = receivedPacket.UnmarshalMsg(res)
	// ch <- p
	//if err1 != nil {
	//	log.Println(err)
	//	conn.Close()
	//}

	fmt.Printf("p.Head.P_Type:%d\n", receivedPacket.Head.P_Type)
	fmt.Printf("p.Head.Op_Type:%d\n", receivedPacket.Head.Op_Type)
	fmt.Printf("p.Payload:%s\n", receivedPacket.Payload)
	SQLContent := string(receivedPacket.Payload)

	err = parser.Parse(strings.NewReader(string(receivedPacket.Payload)), StatementChannel) //收到客户端传过来的语句
	if err != nil {
		log.Fatal(err)
	}
	//close(StatementChannel) //关闭StatementChannel，进而关闭FinishChannel
	//fmt.Println(<-OutputStatementChannel, <-TablesChannel)
	statement := <-OutputStatementChannel
	tabs := <-TablesChannel
	fmt.Println("tabs:", tabs)

	p := CreatePacket(statement, tabs, SQLContent)
	var packetBuf = make([]byte, 0)
	//fmt.Printf("packet.Head.P_Type:%d\n", packet.Head.P_Type)
	//fmt.Printf("packet.Head.Op_Type:%d\n", packet.Head.Op_Type)
	//fmt.Printf("packet.Payload:%s\n", packet.Payload)
	packetBuf, err = p.MarshalMsg(packetBuf)
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
		return
	}
	_, err = conn.Write(packetBuf)
	//err := p.EncodeMsg(wt)
	if err != nil {
		return Type.Packet{}
	}

	return

}

// 监听client的连接，当某个client连接上之后，处理其连接，等待请求的表的信息，获取表名之后将对应的Region地址返回给client
// 如果是创建新表，需要选择一个region，通知region新建表项
// 对于不同的
func main() {
	initMaster()

	StatementChannel = make(chan types.DStatements, 0)
	FinishChannel = make(chan string, 0)
	TablesChannel = make(chan []string, 0)
	OutputStatementChannel = make(chan types.DStatements, 0)
	//FlushChannel := make(chan struct{})
	//fmt.Println("Initialized Master server")
	//var sql_strings = []string{
	//	"create table tab1(a int);",
	//	"select a from tab2;",
	//}

	//err := parser.Parse(strings.NewReader(sql_strings[0]), StatementChannel) //开始解析
	//if err != nil {
	//	log.Fatal(err)
	//}
	go masterDiscovery()
	//regionSer.serverList["127.0.0.1:1234"] = 0
	//regionSer.serverList["127.0.0.2:1234"] = 1
	//regionSer.serverList["127.0.0.3:1234"] = 0
	//CreatePacket(<-OutputStatementChannel, <-TablesChannel, sql_strings[0])
	//close(StatementChannel) //关闭StatementChannel，进而关闭FinishChannel
	//fmt.Println(<-OutputStatementChannel, <-TablesChannel)

	go Parse2Statement(StatementChannel, FinishChannel, TablesChannel, OutputStatementChannel)
	listener, err := net.Listen("tcp", MASTER_PORT)

	for {
		fmt.Println("listening...")
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleClient(conn)
	}
	//for _ = range FinishChannel {
	//
	//}
}

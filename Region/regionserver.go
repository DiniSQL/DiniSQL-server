package Region

import (
	"DiniSQL/MiniSQL"
	"DiniSQL/MiniSQL/src/API"
	"DiniSQL/MiniSQL/src/BufferManager"
	"DiniSQL/MiniSQL/src/CatalogManager"
	"DiniSQL/MiniSQL/src/Interpreter/parser"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"DiniSQL/MiniSQL/src/RecordManager"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

var regionServer *RegionServer
var StatementChannel chan types.DStatements
var FinishChannel chan string
var databaseName string = "minisql"
var FlushChannel chan struct{}

type RegionServer struct {
	tables     []string
	serverID   int
	visitCount int
	register   *ServiceRegister
}

func FlushALl() {
	BufferManager.BlockFlushAll()                                             //缓存block
	RecordManager.FlushFreeList()                                             //free list写回
	CatalogManager.FlushDatabaseMeta(CatalogManager.UsingDatabase.DatabaseId) //刷新记录长度和余量
}

func PeriodicallyFlush() {
	time.Sleep(60 * 5 * time.Second)
	FlushALl()
}

// getLocalIpV4 获取 IPV4 IP，没有则返回空
func getLocalIpV4(interfaceName string) (addr string, err error) {
	inter, err := net.InterfaceByName(interfaceName)
	if err != nil {
		panic(err)
	}
	// 判断网卡是否开启，过滤本地环回接口
	if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
		// 获取网卡下所有的地址
		addrs, err := inter.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				//判断是否存在IPV4 IP 如果没有过滤
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("interface %s don't have an ipv4 address", interfaceName)
}

func InitRegionServer() {
	ip, err := getLocalIpV4("wifi0")
	if err != nil {
		log.Fatal(err)
	}

	var endpoints = []string{"127.0.0.1:2379"}
	prefix := "/region"
	key := strings.Join([]string{prefix, ip}, "/")
	ser, err := NewServiceRegister(endpoints, key, "0", 5)
	if err != nil {
		log.Fatalln(err)
	}

	//监听续租相应chan
	go ser.ListenLeaseRespChan()

	regionServer = new(RegionServer)
	regionServer.visitCount = 0
	regionServer.tables = []string{}
	regionServer.register = ser

	MiniSQL.InitDB()

	StatementChannel = make(chan types.DStatements, 500)   //用于传输操作指令通道
	FinishChannel = make(chan string, 500)                 //用于api执行完成反馈通道
	FlushChannel = make(chan struct{})                     //用于每条指令结束后协程flush
	go API.HandleOneParse(StatementChannel, FinishChannel) //begin the runtime for exec
	go BufferManager.BeginBlockFlush(FlushChannel)
	go PeriodicallyFlush()

	op := "create database " + databaseName + ";"
	err = parser.Parse(strings.NewReader(op), StatementChannel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(<-FinishChannel)
	FlushChannel <- struct{}{} //开始刷新cache

	op = "use database " + databaseName + ";"
	err = parser.Parse(strings.NewReader(op), StatementChannel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(<-FinishChannel)
	FlushChannel <- struct{}{} //开始刷新cache

	// This part is for testing
	op = "create table student2(id int, name char(12) unique, score float,primary key(id) );"
	err = parser.Parse(strings.NewReader(op), StatementChannel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(<-FinishChannel)
	FlushChannel <- struct{}{} //开始刷新cache

	op = "insert into student2 values(1080100001,'name1',99);"
	err = parser.Parse(strings.NewReader(op), StatementChannel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(<-FinishChannel)
	FlushChannel <- struct{}{} //开始刷新cache

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
	// This Upload is for testing
	opRes, err := Upload("student2", "127.0.0.1:3036")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(opRes)
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
	server.register.UpdateKey(fmt.Sprint(server.visitCount))

	var p Packet
	rd := msgp.NewReader(conn)
	wt := msgp.NewWriter(conn)

	err := p.DecodeMsg(rd)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(p.Head.P_Type)

	var opRes string
	if p.Head.P_Type == SQLOperation {
		var statement = byteSliceToString(p.Payload)
		opRes, err = ExecuteSQLOperation(server, statement, p.Head.Op_Type)
		if err != nil {
			fmt.Println(err)
		}
		replyPacket := Packet{Head: PacketHead{P_Type: Result, Op_Type: -1},
			Payload: []byte(opRes)}
		replyPacket.EncodeMsg(wt)
		if err != nil {
			fmt.Println("Error: conn.Write()")
		}
	} else if p.Head.P_Type == UploadRegion {
		str := byteSliceToString(p.Payload)
		strings := strings.Split(str, ",")
		opRes, err = Upload(strings[0], strings[1])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(opRes)
	} else if p.Head.P_Type == RegionTransferPrepare {
		opRes, err, num := PrepareForTransfer(server, p)
		fmt.Println(opRes)
		if err != nil {
			fmt.Println(err)
		}
		for i := 0; i < num; i++ {
			p.DecodeMsg(rd)
			fmt.Println("p.Head.P_Type: ", p.Head.P_Type)
			fmt.Println("p.Head.Op_Type: ", p.Head.Op_Type)
			fmt.Println("p.Head.Spare: ", p.Head.Spare)
			if err != nil {
				fmt.Println("UnmarshalMsg failed")
				fmt.Println(err)
			}
			opRes, err = DownloadTransfer(p)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(opRes)
		}
	} else if p.Head.P_Type == RegionTransfer {
		DownloadTransfer(p)
	}
}

func ExecuteSQLOperation(server *RegionServer, statement string, t types.OperationType) (opRes string, err error) {
	err = parser.Parse(strings.NewReader(string(statement)), StatementChannel)
	if err != nil {
		fmt.Println(err)
	}
	opRes = <-FinishChannel
	FlushChannel <- struct{}{} //开始刷新cache
	fmt.Println(opRes)

	if t == types.CreateTable {
		new_table := strings.Split(opRes, " ")[1]
		server.tables = append(server.tables, new_table)
	} else if t == types.DropTable {
		dropped_table := strings.Split(opRes, " ")[1]
		for i := 0; i < len(server.tables); i++ {
			if server.tables[i] == dropped_table {
				server.tables = append(server.tables[:i], server.tables[i+1:]...)
				break
			}
		}
	}
	return
}

func Upload(table string, to string) (opRes string, err error) {
	catalog := CatalogManager.GetTableCatalogUnsafe(table)
	catalogBuf := make([]byte, 0)
	catalogBuf, err = catalog.MarshalMsg(catalogBuf)
	if err != nil {
		return
	}
	matches, err := filepath.Glob("./data/d_" + databaseName + "_data/" + table + "*")
	if err != nil {
		return
	}
	conn, err := net.Dial("tcp", to)
	if err != nil {
		return
	}

	preparePacket := Packet{Head: PacketHead{P_Type: RegionTransferPrepare, Op_Type: len(matches), Spare: table},
		Payload: catalogBuf}
	wt := msgp.NewWriter(conn)
	err = preparePacket.EncodeMsg(wt)
	if err != nil {
		return
	}

	for i := range matches {
		fmt.Println(matches[i])
		file, err := os.Open(matches[i])
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		packet := Packet{Head: PacketHead{P_Type: RegionTransfer, Op_Type: i, Spare: strings.Split(matches[i], "\\")[2]},
			Payload: content}
		err = packet.EncodeMsg(wt)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(3 * time.Second)
	}
	opRes = "Upload Success"
	return
}

func PrepareForTransfer(server *RegionServer, p Packet) (opRes string, err error, num int) {
	catalog := CatalogManager.TableCatalog{}
	_, err = catalog.UnmarshalMsg(p.Payload)
	if err != nil {
		return
	}
	fmt.Println(string(p.Payload))
	server.tables = append(server.tables, p.Head.Spare)
	CatalogManager.TableName2CatalogMap[p.Head.Spare] = &catalog
	err = CatalogManager.FlushDatabaseMeta(CatalogManager.UsingDatabase.DatabaseId)
	if err != nil {
		return
	}
	num = p.Head.Op_Type
	opRes = "Prepare Success, ready to download file"
	return
}

func DownloadTransfer(p Packet) (opRes string, err error) {
	filename := p.Head.Spare
	dir := "./data/d_" + databaseName + "_data/" + filename
	file, err := os.Create(dir)
	if err != nil {
		return
	}
	_, err = file.Write(p.Payload)
	if err == nil {
		opRes = filename + " download success"
	}
	return
}

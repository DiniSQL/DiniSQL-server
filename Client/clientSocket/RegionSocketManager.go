package Socket

import (
    "fmt"
    "log"
    "net"
	"DiniSQL/Client/type"
	"github.com/tinylib/msgp/msgp"
)

var ch = make(chan []byte, 5)

// initial connection to region server and send message
func ConnectToRegion(regionIP string,regionPort int, packet Type.Packet)bool{

    address := net.TCPAddr{
        IP:   net.ParseIP(regionIP), // 把字符串IP地址转换为net.IP类型
        Port: regionPort,
    }
	conn, err := net.DialTCP("tcp4", nil, &address)
    if err != nil {
        log.Fatal(err) // Println + os.Exit(1)
		return false
    }
	var packetBuf = make([]byte, packet.Msgsize())
	packetBuf,err = packet.MarshalMsg(packetBuf)
	fmt.Println(packetBuf)
	if err != nil {
        log.Fatal(err) // Println + os.Exit(1)
		return false
    }
	// packetBuf[0] = 3
	// packetBuf[1] = 3
	// packetBuf[2] = 3

	_, err1 := conn.Write(packetBuf)
	fmt.Println(packet)
	if err1 != nil {
		log.Println(err)
		conn.Close()
		return false
	}
	fmt.Printf("send %d to %s\n", packet.Head.P_Type, conn.RemoteAddr())
	conn.Close()

	return true

}

// listen 
// input : IP and Port of client
func KeepListening(ClientIP string, ClientPort int){
    address := net.TCPAddr{
        IP:   net.ParseIP(ClientIP), // 把字符串IP地址转换为net.IP类型
        Port: ClientPort,
    }
    listener, err := net.ListenTCP("tcp4", &address) // 创建TCP4服务器端监听器
    if err != nil {
        log.Fatal(err) // Println + os.Exit(1)
    }
    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            log.Fatal(err) // 错误直接退出
        }
        fmt.Println("remote address:", conn.RemoteAddr())
        res := make([]byte,1024)
        _, err1 := conn.Read(res)
		ch <- res
		if err1 != nil {
            log.Println(err)
            conn.Close()
        }
        fmt.Printf("received %s from %s", res, conn.RemoteAddr())
    }
}

// func main(){
// 	var sql string
// 	go KeepListening("127.0.0.1",8004)
// 	for true  {
//         fmt.Scanln(&sql)
// 		p := Region.Packet{Head: Region.PacketHead{P_Type: Region.KeepAlive, Op_Type: -1},
// 			Payload: []byte(sql)}
// 		ConnectToRegion("10.192.16.210",6100,p)
//     }
// }
package Master

import "net"

type RegionStatus struct {
	regionID  int
	surviving bool
	loadLevel int
}

const (
	MASTER_PORT = ":6100"
)

var tab2reg map[string]int
var id2reg map[int]*RegionStatus

func initMaster() {
	// Initialize the master server
	tab2reg = make(map[string]int) // table name to region id
	tab2reg["foo"] = 1
	id2reg = make(map[int]*RegionStatus) // region id to region status
	id2reg[0] = new(RegionStatus)
	id2reg[0].regionID = 1

}

func listenMaster() {
	initMaster()
	listen, err := net.Listen("tcp", MASTER_PORT)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go handleMaster(conn)
	}
}

func handleMaster(conn net.Conn) {
	defer conn.Close()
	for {
		// handle request from client
	}
}

package Region

type PacketType = int

const (
	KeepAlive    PacketType = iota // RegionServer send to Master
	Ask                            // Client send to Master to know which RegionServer it should visit
	Answer                         // Master answer the Ask packet from Client, tell which RegionServer
	SQLOperation                   // Client send to RegionServer to execute a SQL operation
	Result                         // RegionServer send to Client, the result of the SQL operation
	RegionBackup                   // Master send to RegionServer, to build a backup of a region
)

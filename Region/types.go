package Region

type PacketType = int

const (
	KeepAlive             PacketType = iota // RegionServer send to Master
	Ask                                     // Client send to Master to know which RegionServer it should visit
	Answer                                  // Master answer the Ask packet from Client, tell which RegionServer
	SQLOperation                            // Client send to RegionServer to execute a SQL operation
	Result                                  // RegionServer send to Client, the result of the SQL operation
	UploadRegion                            // Master send to RegionServer, tell the RegionServer to upload a region to another server
	RegionTransferPrepare                   // RegionServer send to RegionServer, transfer catalog before transfering file
	RegionTransfer                          // RegionServer send to RegionServer, transfer all files of a table
)

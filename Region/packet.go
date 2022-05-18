package Region

// import (
//     "fmt"
// )

//go:generate msgp

type PacketHead struct {
	P_Type  int
	Op_Type int
	Spare   string
}

type Packet struct {
	Head    PacketHead
	Payload []byte
}

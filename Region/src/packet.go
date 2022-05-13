package Region

import (
	"DiniSQL/MiniSQL/src/Interpreter/types"
)

//go:generate msgp

type PacketHead struct {
	packetType   PacketType
	detailedType types.OperationType
}

type Payload struct {
	content []byte
}

type Packet struct {
	head    PacketHead
	payload Payload
}

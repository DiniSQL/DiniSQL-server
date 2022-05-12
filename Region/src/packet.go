package Region

import (
	"DiniSQL/MiniSQL/src/Interpreter/types"
)

//go:generate msgp

type PacketType = int

const (
	KeepAlive PacketType = iota
	Ask
	Answer
	SQLOperation
	Result
)

type PacketHead struct {
	packetType   PacketType
	detailedType types.OperationType
	payload      Payload
}

type Payload struct {
	content []byte
}

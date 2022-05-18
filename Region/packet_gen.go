package Region

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Packet) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Head":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "Head")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					err = msgp.WrapError(err, "Head")
					return
				}
				switch msgp.UnsafeString(field) {
				case "P_Type":
					z.Head.P_Type, err = dc.ReadInt()
					if err != nil {
						err = msgp.WrapError(err, "Head", "P_Type")
						return
					}
				case "Op_Type":
					z.Head.Op_Type, err = dc.ReadInt()
					if err != nil {
						err = msgp.WrapError(err, "Head", "Op_Type")
						return
					}
				case "Spare":
					z.Head.Spare, err = dc.ReadString()
					if err != nil {
						err = msgp.WrapError(err, "Head", "Spare")
						return
					}
				default:
					err = dc.Skip()
					if err != nil {
						err = msgp.WrapError(err, "Head")
						return
					}
				}
			}
		case "Payload":
			z.Payload, err = dc.ReadBytes(z.Payload)
			if err != nil {
				err = msgp.WrapError(err, "Payload")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Packet) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Head"
	err = en.Append(0x82, 0xa4, 0x48, 0x65, 0x61, 0x64)
	if err != nil {
		return
	}
	// map header, size 3
	// write "P_Type"
	err = en.Append(0x83, 0xa6, 0x50, 0x5f, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Head.P_Type)
	if err != nil {
		err = msgp.WrapError(err, "Head", "P_Type")
		return
	}
	// write "Op_Type"
	err = en.Append(0xa7, 0x4f, 0x70, 0x5f, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Head.Op_Type)
	if err != nil {
		err = msgp.WrapError(err, "Head", "Op_Type")
		return
	}
	// write "Spare"
	err = en.Append(0xa5, 0x53, 0x70, 0x61, 0x72, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Head.Spare)
	if err != nil {
		err = msgp.WrapError(err, "Head", "Spare")
		return
	}
	// write "Payload"
	err = en.Append(0xa7, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Payload)
	if err != nil {
		err = msgp.WrapError(err, "Payload")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Packet) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Head"
	o = append(o, 0x82, 0xa4, 0x48, 0x65, 0x61, 0x64)
	// map header, size 3
	// string "P_Type"
	o = append(o, 0x83, 0xa6, 0x50, 0x5f, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.Head.P_Type)
	// string "Op_Type"
	o = append(o, 0xa7, 0x4f, 0x70, 0x5f, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.Head.Op_Type)
	// string "Spare"
	o = append(o, 0xa5, 0x53, 0x70, 0x61, 0x72, 0x65)
	o = msgp.AppendString(o, z.Head.Spare)
	// string "Payload"
	o = append(o, 0xa7, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64)
	o = msgp.AppendBytes(o, z.Payload)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Packet) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Head":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Head")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					err = msgp.WrapError(err, "Head")
					return
				}
				switch msgp.UnsafeString(field) {
				case "P_Type":
					z.Head.P_Type, bts, err = msgp.ReadIntBytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "Head", "P_Type")
						return
					}
				case "Op_Type":
					z.Head.Op_Type, bts, err = msgp.ReadIntBytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "Head", "Op_Type")
						return
					}
				case "Spare":
					z.Head.Spare, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "Head", "Spare")
						return
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						err = msgp.WrapError(err, "Head")
						return
					}
				}
			}
		case "Payload":
			z.Payload, bts, err = msgp.ReadBytesBytes(bts, z.Payload)
			if err != nil {
				err = msgp.WrapError(err, "Payload")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Packet) Msgsize() (s int) {
	s = 1 + 5 + 1 + 7 + msgp.IntSize + 8 + msgp.IntSize + 6 + msgp.StringPrefixSize + len(z.Head.Spare) + 8 + msgp.BytesPrefixSize + len(z.Payload)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PacketHead) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "P_Type":
			z.P_Type, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "P_Type")
				return
			}
		case "Op_Type":
			z.Op_Type, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Op_Type")
				return
			}
		case "Spare":
			z.Spare, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Spare")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z PacketHead) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "P_Type"
	err = en.Append(0x83, 0xa6, 0x50, 0x5f, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.P_Type)
	if err != nil {
		err = msgp.WrapError(err, "P_Type")
		return
	}
	// write "Op_Type"
	err = en.Append(0xa7, 0x4f, 0x70, 0x5f, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Op_Type)
	if err != nil {
		err = msgp.WrapError(err, "Op_Type")
		return
	}
	// write "Spare"
	err = en.Append(0xa5, 0x53, 0x70, 0x61, 0x72, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Spare)
	if err != nil {
		err = msgp.WrapError(err, "Spare")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z PacketHead) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "P_Type"
	o = append(o, 0x83, 0xa6, 0x50, 0x5f, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.P_Type)
	// string "Op_Type"
	o = append(o, 0xa7, 0x4f, 0x70, 0x5f, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.Op_Type)
	// string "Spare"
	o = append(o, 0xa5, 0x53, 0x70, 0x61, 0x72, 0x65)
	o = msgp.AppendString(o, z.Spare)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PacketHead) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "P_Type":
			z.P_Type, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "P_Type")
				return
			}
		case "Op_Type":
			z.Op_Type, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Op_Type")
				return
			}
		case "Spare":
			z.Spare, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Spare")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z PacketHead) Msgsize() (s int) {
	s = 1 + 7 + msgp.IntSize + 8 + msgp.IntSize + 6 + msgp.StringPrefixSize + len(z.Spare)
	return
}

package connector

import "encoding/binary"

type CmdType uint16

const (
	CmdTypeUnknown CmdType = iota
	CmdTypeHeartbeat
)

const (
	PacketHeaderSize = (16 + 16 + 32) / 8
	PacketVersion    = 1
)

type Header struct {
	// Packet protocol version
	version uint16
	// Packet command
	cmd CmdType
	// size of payload
	size uint32
}

type Packet struct {
	header  *Header
	payload []byte
}

func (p *Packet) ParseHeader(header []byte) error {
	if len(header) != PacketHeaderSize {
		return ErrPacketHeaderInvalid
	}

	p.header = &Header{}
	p.header.version = binary.BigEndian.Uint16(header[0:2])
	p.header.cmd = CmdType(binary.BigEndian.Uint16(header[2:4]))
	p.header.size = binary.BigEndian.Uint32(header[4:8])

	if p.header.version != PacketVersion {
		return ErrPacketVersion
	}
	if p.header.cmd <= 0 {
		return ErrPacketCommand
	}
	return nil
}

func NewEmptyPacket() *Packet {
	return &Packet{
		header:  &Header{},
		payload: nil,
	}
}

func NewPacket(header *Header, payload []byte) *Packet {
	return &Packet{
		header:  header,
		payload: payload,
	}
}

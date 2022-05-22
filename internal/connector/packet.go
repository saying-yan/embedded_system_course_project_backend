package connector

import (
	"encoding/binary"
	"fmt"
)

type CmdType uint16

const (
	CmdTypeUnknown CmdType = iota
	CmdTypeHeartbeat
	CmdTypeDeviceInfo
	CmdTypeSongsInfo
	CmdTypeExit
	CmdTypeMusicCompleted

	CmdTypePlayMusic CmdType = iota + 1000
)

func (c CmdType) String() string {
	switch c {
	case CmdTypeHeartbeat:
		return "CmdTypeHeartbeat"
	case CmdTypeDeviceInfo:
		return "CmdTypeDeviceInfo"
	case CmdTypeSongsInfo:
		return "CmdTypeSongsInfo"
	default:
		return "CmdTypeUnknown"
	}
}

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

func (p *Packet) String() string {
	return fmt.Sprintf("%s packet with paylaod: %v", p.header.cmd.String(), p.payload)
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

func (p *Packet) Bytes() []byte {
	buf := make([]byte, PacketHeaderSize, PacketHeaderSize+p.header.size)
	binary.BigEndian.PutUint16(buf, p.header.version)
	binary.BigEndian.PutUint16(buf[2:], uint16(p.header.cmd))
	binary.BigEndian.PutUint32(buf[4:], p.header.size)
	if len(p.payload) > 0 {
		buf = append(buf, p.payload...)
	}
	return buf
}

func NewEmptyPacket() *Packet {
	return &Packet{
		header: &Header{
			version: PacketVersion,
		},
		payload: nil,
	}
}

func NewPacket(header *Header, payload []byte) *Packet {
	return &Packet{
		header:  header,
		payload: payload,
	}
}

func (p *Packet) WithCmd(cmd CmdType) *Packet {
	p.header.cmd = cmd
	return p
}

func (p *Packet) WithPayload(payload []byte) *Packet {
	p.header.size = uint32(len(payload))
	p.payload = payload
	return p
}

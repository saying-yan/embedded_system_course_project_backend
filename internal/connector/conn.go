package connector

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync/atomic"
	"time"

	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
)

type ConnState int

const (
	ConnStateUnknown ConnState = iota
	ConnStateActive
	ConnStateIdle
	ConnStateClosed
)

var connIDCounter uint64 = 1000

type Conn struct {
	ID         uint64
	RemoteAddr string

	netConn    net.Conn
	activeTime time.Time
	exitChan   chan struct{}
}

func (conn *Conn) String() string {
	return fmt.Sprintf("connection %d from %s, activeTime: %s", conn.ID, conn.RemoteAddr, conn.activeTime.String())
}

func (conn *Conn) receivePacket() (*Packet, error) {
	headerBuf, err := ioutil.ReadAll(io.LimitReader(conn.netConn, PacketHeaderSize))
	if err != nil {
		return nil, err
	}
	Logger.Debugf("receive header: %v", headerBuf)

	if len(headerBuf) <= 0 {
		return nil, io.EOF
	}
	conn.activeTime = time.Now()
	packet := NewEmptyPacket()
	err = packet.ParseHeader(headerBuf)
	if err != nil {
		return nil, err
	}

	payloadSize := packet.header.size
	if payloadSize == 0 {
		return packet, nil
	}

	packet.payload, err = ioutil.ReadAll(io.LimitReader(conn.netConn, int64(payloadSize)))
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func (conn *Conn) handleConn() {
	// handle received packet
	for {
		packet, err := conn.receivePacket()
		if err != nil {
			Logger.Errorf("conn:%d from %s receive packet error: %s", conn.ID, conn.RemoteAddr, err.Error())
			continue
		}

		handler := commandHandlerMap[packet.header.cmd]
		go handler(conn, packet)
	}
}

func (conn *Conn) Close() {
	// TODO: close conn
}

func NewConn(rawConn net.Conn) *Conn {
	return &Conn{
		ID:         atomic.AddUint64(&connIDCounter, 1),
		RemoteAddr: rawConn.RemoteAddr().String(),
		netConn:    rawConn,
		activeTime: time.Now(),
		exitChan:   make(chan struct{}),
	}
}

package connector

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync/atomic"
	"time"

	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
)

//type ConnState int
//
//const (
//	ConnStateUnknown ConnState = iota
//	ConnStateActive
//	ConnStateIdle
//	ConnStateClosed
//)

type Conn struct {
	DeviceID   atomic.Value
	RemoteAddr string

	netConn    net.Conn
	activeTime atomic.Value
	exitChan   chan struct{}
	exited     int32
}

func (conn *Conn) String() string {
	return fmt.Sprintf("connection %d from %s, activeTime: %s", conn.getDeviceID(), conn.RemoteAddr, conn.getActiveTime().String())
}

func (conn *Conn) getActiveTime() time.Time {
	t, ok := conn.activeTime.Load().(time.Time)
	if !ok {
		// impossible
		Logger.Errorf("get activeTime error")
		return time.Time{}
	}
	return t
}

func (conn *Conn) setActiveTime(t time.Time) {
	conn.activeTime.Store(t)
}

func (conn *Conn) setDeviceID(deviceID uint32) {
	_, ok := conn.DeviceID.Load().(uint32)
	if ok {
		Logger.Errorf("already set deviceID")
		panic("already set deviceID")
	}
	conn.DeviceID.Store(deviceID)

	// put into connPool
	ConnPool.PutConn(conn)
}

func (conn *Conn) getDeviceID() uint32 {
	id, ok := conn.DeviceID.Load().(uint32)
	if !ok {
		// impossible
		Logger.Errorf("get DeviceID error")
		return 0
	}
	return id
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
	conn.setActiveTime(time.Now())
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
		if atomic.LoadInt32(&conn.exited) > 0 {
			break
		}

		packet, err := conn.receivePacket()
		if err != nil {
			Logger.Debugf("conn:%d from %s receive packet error: %s", conn.getDeviceID(), conn.RemoteAddr, err.Error())
			if err == io.EOF {
				conn.Close()
				break
			}
			continue
		}

		handler := commandHandlerMap[packet.header.cmd]
		go func() {
			//defer func() {
			//	rc := recover()
			//	if rc != nil {
			//		// TODO: recover
			//	}
			//}()
			if handler == nil {
				Logger.Errorf("handle conn:%d from %s, packet:%s error: unknown cmd", conn.getDeviceID(), conn.RemoteAddr, packet.String())
				return
			}
			err := handler(conn, packet)
			if err != nil {
				Logger.Errorf("handle conn:%d from %s, packet:%s error: %s", conn.getDeviceID(), conn.RemoteAddr, packet.String(), err.Error())
			}
			conn.setActiveTime(time.Now())
		}()
	}
}

func (conn *Conn) PlayMusic(songID uint32) error {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, songID)
	packet := NewEmptyPacket().WithCmd(CmdTypePlayMusic).WithPayload(payload)

	buf := packet.Bytes()
	_, err := conn.netConn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Conn) Close() {
	// close raw conn
	conn.netConn.Close()
	// close channel
	select {
	case <-conn.exitChan:
	default:
		close(conn.exitChan)
	}
	delete(ConnPool.connMap, conn.getDeviceID())
	atomic.StoreInt32(&conn.exited, 1)
}

func NewConn(rawConn net.Conn) *Conn {
	conn := &Conn{
		RemoteAddr: rawConn.RemoteAddr().String(),
		netConn:    rawConn,
		exitChan:   make(chan struct{}),
	}
	conn.setActiveTime(time.Now())
	return conn
}

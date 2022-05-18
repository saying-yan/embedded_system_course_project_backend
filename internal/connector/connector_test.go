package connector

import (
	"encoding/binary"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"net"
	"testing"
	"time"
)

func initLogger() {
	InitLogger("debug", "../../log/backend", true)
}

func TestConnector(t *testing.T) {
	initLogger()
	c, err := NewConnector(8888)
	if err != nil {
		t.Fatalf("new connector error: %s", err.Error())
	}

	go func() {
		err = c.Serve()
		if err != nil {
			Logger.Fatalf("new connector error: %s", err.Error())
		}
	}()

	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatalf("dial tcp error: %s", err.Error())
	}

	deviceID := make([]byte, 8)
	binary.BigEndian.PutUint64(deviceID, 1000)
	packet := NewPacket(&Header{
		version: PacketVersion,
		cmd:     CmdTypeDeviceInfo,
		size:    8,
	}, deviceID)
	buf := packet.Bytes()

	_, err = conn.Write(buf)
	time.Sleep(300 * time.Millisecond)

	packet = NewPacket(&Header{
		version: PacketVersion,
		cmd:     CmdTypeHeartbeat,
		size:    0,
	}, nil)
	buf = packet.Bytes()

	_, err = conn.Write(buf)

	if err != nil {
		t.Fatalf("write tcp error: %s", err.Error())
	}

	packet = NewEmptyPacket().WithCmd(CmdTypeHeartbeat)
	buf = packet.Bytes()
	_, err = conn.Write(buf)

	packet = NewEmptyPacket().WithCmd(CmdTypeExit)
	buf = packet.Bytes()
	_, err = conn.Write(buf)

	if err != nil {
		t.Fatalf("write tcp error: %s", err.Error())
	}

	time.Sleep(8 * time.Second)
	for _, serverConn := range connPool.connMap {
		Logger.Debugf("server connection: %s", serverConn.String())
	}
}

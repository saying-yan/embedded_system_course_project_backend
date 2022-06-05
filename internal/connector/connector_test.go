package connector

import (
	"encoding/binary"
	"fmt"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/provider"
	"net"
	"testing"
	"time"
)

func initLogger() {
	InitLogger("debug", "../../log/backend", true)
}

func TestConnector(t *testing.T) {
	initLogger()
	c, err := NewConnector(8002, 5*time.Second)
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
	conn, err := net.Dial("tcp", "saying.mobi:8002")
	if err != nil {
		t.Fatalf("dial tcp error: %s", err.Error())
	}

	deviceID := make([]byte, 4)
	binary.BigEndian.PutUint32(deviceID, 1000)
	packet := NewPacket(&Header{
		version: PacketVersion,
		cmd:     CmdTypeDeviceInfo,
		size:    4,
	}, deviceID)
	buf := packet.Bytes()

	fmt.Println("send buff", buf)
	_, err = conn.Write(buf)
	time.Sleep(300 * time.Millisecond)

	packet = NewEmptyPacket().WithCmd(CmdTypeHeartbeat)
	buf = packet.Bytes()
	fmt.Println("send buff", buf)
	_, err = conn.Write(buf)

	if err != nil {
		t.Fatalf("write tcp error: %s", err.Error())
	}

	time.Sleep(2 * time.Second)
	packet = NewEmptyPacket().WithCmd(CmdTypeHeartbeat)
	buf = packet.Bytes()
	_, err = conn.Write(buf)

	//packet = NewEmptyPacket().WithCmd(CmdTypeExit)
	//buf = packet.Bytes()
	//_, err = conn.Write(buf)

	packet = NewEmptyPacket().WithCmd(CmdTypeSongsInfo)
	songInfo := []byte("\x00\x00\x00\x01\x00\x04\x00\x06namesinger\x00\x00\x00\x02\x00\x05\x00\x07name2singer2")
	packet = packet.WithPayload(songInfo)
	buf = packet.Bytes()
	_, err = conn.Write(buf)

	time.Sleep(1 * time.Second)
	for _, serverConn := range ConnPool.connMap {
		Logger.Debugf("server connection: %s", serverConn.String())
	}
	for _, device := range provider.Provider.Devices {
		for _, song := range device.Songs {
			Logger.Debugf("song: %#v", song)
		}
	}
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			packet = NewEmptyPacket().WithCmd(CmdTypeHeartbeat)
			_, err = conn.Write(packet.Bytes())
			fmt.Println("heartbeat")
		}
	}
}

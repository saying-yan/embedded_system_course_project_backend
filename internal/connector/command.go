package connector

import (
	"encoding/binary"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/provider"
)

type commandHandler func(conn *Conn, packet *Packet) error

var commandHandlerMap = map[CmdType]commandHandler{
	CmdTypeUnknown:    UnknownHandler,
	CmdTypeHeartbeat:  HeartbeatHandler,
	CmdTypeDeviceInfo: DeviceInfoHandler,
	CmdTypeSongsInfo:  SongsInfoHandler,
	CmdTypeExit:       ExitHandler,
}

func UnknownHandler(_ *Conn, _ *Packet) error {
	panic("unknown command type in packet")
}

func HeartbeatHandler(_ *Conn, _ *Packet) error {
	return nil
}

func DeviceInfoHandler(conn *Conn, p *Packet) error {
	size := p.header.size
	payload := p.payload

	if size != 8 {
		Logger.Debugf("DeviceInfo packet error: patload (size=%d) != 8", size)
		return ErrPacketDeviceInfo
	}

	deviceID := binary.BigEndian.Uint64(payload)
	conn.setDeviceID(deviceID)

	deviceInfo := provider.NewDeviceInfo(deviceID, conn.RemoteAddr)
	provider.Provider.SetDeviceInfo(deviceInfo)

	return nil
}

func SongsInfoHandler(conn *Conn, p *Packet) error {
	size := p.header.size
	payload := p.payload

	var sum uint32 = 0

	var songs []*provider.Song
	for sum < size {
		songID := binary.BigEndian.Uint32(payload[sum : sum+4])
		nameLen := binary.BigEndian.Uint16(payload[sum+4 : sum+6])
		tmp := sum + 6 + uint32(nameLen)
		name := string(payload[sum+6 : tmp])
		singerNameLen := binary.BigEndian.Uint16(payload[tmp : tmp+2])
		singerName := string(payload[tmp+2 : tmp+2+uint32(singerNameLen)])
		songInfo := provider.NewSong(songID, name, singerName)
		songs = append(songs, songInfo)
		sum = tmp + 2 + uint32(singerNameLen)
	}
	return provider.Provider.AddSongs(conn.getDeviceID(), songs)
}

func ExitHandler(conn *Conn, p *Packet) error {
	connPool.removeConn(conn.getDeviceID())
	conn.Close()
	return nil
}

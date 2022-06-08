package connector

import (
	"encoding/binary"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/provider"
)

type commandHandler func(conn *Conn, packet *Packet) error

var commandHandlerMap = map[CmdType]commandHandler{
	CmdTypeUnknown:        UnknownHandler,
	CmdTypeHeartbeat:      HeartbeatHandler,
	CmdTypeDeviceInfo:     DeviceInfoHandler,
	CmdTypeSongsInfo:      SongsInfoHandler,
	CmdTypeExit:           ExitHandler,
	CmdTypeMusicCompleted: MusicCompletedHandler,
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

	if size != 4 {
		Logger.Debugf("DeviceInfo packet error: payload (size=%d) != 4", size)
		return ErrPacketDeviceInfo
	}

	deviceID := binary.BigEndian.Uint32(payload)
	conn.setDeviceID(deviceID)

	deviceInfo := provider.NewDeviceInfo(deviceID, conn.RemoteAddr)
	provider.SetDeviceInfo(deviceInfo)

	return nil
}

func SongsInfoHandler(conn *Conn, p *Packet) error {
	size := p.header.size
	payload := p.payload

	var index uint32 = 0

	var songs []*provider.Song
	for index < size {
		if index+8 > size {
			return ErrPacketSongInfo
		}
		songID := binary.BigEndian.Uint32(payload[index : index+4])
		nameLen := uint32(binary.BigEndian.Uint16(payload[index+4 : index+6]))
		singerNameLen := uint32(binary.BigEndian.Uint16(payload[index+6 : index+8]))
		tmp := index + 8
		if tmp+nameLen+singerNameLen > size {
			return ErrPacketSongInfo
		}

		name := string(payload[tmp : tmp+nameLen])
		singerName := string(payload[tmp+nameLen : tmp+nameLen+singerNameLen])

		songInfo := provider.NewSong(songID, name, singerName)
		Logger.Debugf("get song info: %#v", songInfo)
		songs = append(songs, songInfo)
		index = tmp + nameLen + singerNameLen
	}
	return provider.GetDeviceProvider(conn.getDeviceID()).AddSongs(songs)
}

func MusicCompletedHandler(conn *Conn, _ *Packet) error {
	Logger.Debugf("[conn %d] receive music completed", conn.getDeviceID())
	nextSongID := provider.GetDeviceProvider(conn.getDeviceID()).GetNextSongID()
	Logger.Debugf("[conn %d] next music %d", conn.getDeviceID(), nextSongID)
	return conn.PlayMusic(nextSongID)
}

func ExitHandler(conn *Conn, p *Packet) error {
	conn.Close()
	return nil
}

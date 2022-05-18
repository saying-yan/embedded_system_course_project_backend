package provider

// DeviceInfo 嵌入式设备信息
type DeviceInfo struct {
	ID   uint64
	Addr string
}

// Song 歌曲信息
type Song struct {
	SongID     uint32 // 歌曲ID
	Name       string // 歌曲名
	SingerName string // 歌手名
	// Duration   string // 乐曲时长
}

func NewSong(id uint32, name, singerName string) *Song {
	return &Song{
		SongID:     id,
		Name:       name,
		SingerName: singerName,
	}
}

func NewDeviceInfo(id uint64, addr string) *DeviceInfo {
	return &DeviceInfo{
		ID:   id,
		Addr: addr,
	}
}

//func (d *DeviceInfo) addSongWithoutLock(song *Song) error {
//	_, exists := d.Songs[song.SongID]
//	if exists {
//		return ErrSongIDAlreadyExists
//	}
//	d.Songs[song.SongID] = song
//	return nil
//}

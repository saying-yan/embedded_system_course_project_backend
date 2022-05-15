package provider

// DeviceInfo 嵌入式设备信息
type DeviceInfo struct {
	ID    int64
	Addr  string
	Songs map[int64]Song
}

// Song 歌曲信息
type Song struct {
	SongID     int64  // 歌曲ID
	Name       string // 歌曲名
	SingerName string // 歌手名
	Duration   string // 乐曲时长
}

func NewDeviceInfo(id int64, addr string) *DeviceInfo {
	return &DeviceInfo{
		ID:    id,
		Addr:  addr,
		Songs: make(map[int64]Song),
	}
}

func (d *DeviceInfo) AddSong(song Song) error {
	_, exists := d.Songs[song.SongID]
	if exists {
		return ErrSongIDAlreadyExists
	}
	d.Songs[song.SongID] = song
	return nil
}

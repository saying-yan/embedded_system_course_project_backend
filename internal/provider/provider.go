package provider

import "sync"

var Provider = newMemoryProvider()

type MemoryProvider struct {
	Devices map[uint32]*DeviceProvider
}

const (
	TotalList = iota
	OrderedList
)

// DeviceProvider 为了简化实现，直接将数据存入内存
type DeviceProvider struct {
	DeviceInfo  *DeviceInfo
	Songs       map[uint32]*Song
	TotalList   []uint32
	OrderedList []uint32
	rwLock      sync.RWMutex
}

func GetDeviceProvider(deviceID uint32) *DeviceProvider {
	return Provider.Devices[deviceID]
}

func SetDeviceInfo(d *DeviceInfo) {
	// should set only once, no need to lock
	device, ok := Provider.Devices[d.ID]
	if !ok {
		Provider.Devices[d.ID] = &DeviceProvider{
			DeviceInfo:  d,
			Songs:       make(map[uint32]*Song),
			TotalList:   nil,
			OrderedList: nil,
		}
		device = Provider.Devices[d.ID]
	}
	device.DeviceInfo = d
}

func (p *DeviceProvider) AddSongs(songs []*Song) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	for _, s := range songs {
		_, ok := p.Songs[s.SongID]
		if ok {
			// already exists, only update song info
			p.Songs[s.SongID] = s
			continue
		}

		p.Songs[s.SongID] = s
		p.TotalList = append(p.TotalList, s.SongID)
	}
	return nil
}

func (p *DeviceProvider) GetList(listType int) ([]*Song, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()

	var songs []*Song
	switch listType {
	case TotalList:
		songs = make([]*Song, 0, len(p.TotalList))
		for _, id := range p.TotalList {
			songs = append(songs, p.Songs[id])
		}
	case OrderedList:
		songs = make([]*Song, 0, len(p.OrderedList))
		for _, id := range p.OrderedList {
			songs = append(songs, p.Songs[id])
		}
	default:
		return nil, ErrUnknownListType
	}
	return songs, nil
}

func (p *DeviceProvider) OrderSong(songID uint32) error {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()

	if _, ok := p.Songs[songID]; !ok {
		return ErrSongNotExists
	}

	p.OrderedList = append(p.OrderedList, songID)
	return nil
}

func newMemoryProvider() *MemoryProvider {
	return &MemoryProvider{Devices: make(map[uint32]*DeviceProvider)}
}

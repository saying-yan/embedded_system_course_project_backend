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
	totalList   []uint32
	orderedList []uint32
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
			totalList:   nil,
			orderedList: nil,
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
		p.totalList = append(p.totalList, s.SongID)
	}
	return nil
}

func (p *DeviceProvider) GetNextSongID() uint32 {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	if len(p.orderedList) == 0 {
		return 0
	} else if len(p.orderedList) == 1 {
		songID := p.orderedList[0]
		p.orderedList = nil
		return songID
	} else {
		songID := p.orderedList[0]
		p.orderedList = p.orderedList[1:]
		return songID
	}
}

func (p *DeviceProvider) GetList(listType int) ([]*Song, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()

	var songs []*Song
	switch listType {
	case TotalList:
		songs = make([]*Song, 0, len(p.totalList))
		for _, id := range p.totalList {
			songs = append(songs, p.Songs[id])
		}
	case OrderedList:
		songs = make([]*Song, 0, len(p.orderedList))
		for _, id := range p.orderedList {
			songs = append(songs, p.Songs[id])
		}
	default:
		return nil, ErrUnknownListType
	}
	return songs, nil
}

func (p *DeviceProvider) OrderSong(songID uint32) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	if _, ok := p.Songs[songID]; !ok {
		return ErrSongNotExists
	}

	p.orderedList = append(p.orderedList, songID)
	return nil
}

func (p *DeviceProvider) StickTopSong(songIndex int) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	if songIndex >= len(p.orderedList) {
		return ErrSongNotExists
	}
	if songIndex == 0 {
		return nil
	}

	temp := p.orderedList[songIndex]
	for i := songIndex; i > 0; i-- {
		p.orderedList[songIndex] = p.orderedList[songIndex-1]
	}
	p.orderedList[0] = temp
	return nil
}

func newMemoryProvider() *MemoryProvider {
	return &MemoryProvider{Devices: make(map[uint32]*DeviceProvider)}
}

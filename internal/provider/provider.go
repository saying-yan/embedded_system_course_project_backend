package provider

import "sync"

var Provider = newMemoryProvider()

type MemoryProvider struct {
	Devices map[uint64]*DeviceProvider
}

// DeviceProvider 为了简化实现，直接将数据存入内存
type DeviceProvider struct {
	DeviceInfo  *DeviceInfo
	Songs       map[uint32]*Song
	TotalList   []uint32
	OrderedList []uint32
	rwLock      sync.RWMutex
}

func (p *MemoryProvider) SetDeviceInfo(d *DeviceInfo) {
	// should set only once, no need to lock
	device, ok := p.Devices[d.ID]
	if !ok {
		p.Devices[d.ID] = &DeviceProvider{
			DeviceInfo:  d,
			Songs:       make(map[uint32]*Song),
			TotalList:   nil,
			OrderedList: nil,
		}
		device = p.Devices[d.ID]
	}
	device.DeviceInfo = d
}

func (p *MemoryProvider) AddSongs(deviceID uint64, songs []*Song) error {
	device, ok := p.Devices[deviceID]
	if !ok {
		return ErrDeviceNotExists
	}

	device.rwLock.Lock()
	defer device.rwLock.Unlock()

	for _, s := range songs {
		_, ok = device.Songs[s.SongID]
		if ok {
			// already exists, only update song info
			device.Songs[s.SongID] = s
			continue
		}

		device.Songs[s.SongID] = s
		device.TotalList = append(device.TotalList, s.SongID)
	}
	return nil
}

func (p *MemoryProvider) GetTotalList(deviceID uint64) ([]*Song, error) {
	device, ok := p.Devices[deviceID]
	if !ok {
		return nil, ErrDeviceNotExists
	}
	device.rwLock.RLock()
	defer device.rwLock.RUnlock()

	songs := make([]*Song, 0, len(device.TotalList))
	for _, id := range device.TotalList {
		songs = append(songs, device.Songs[id])
	}
	return songs, nil
}

func newMemoryProvider() *MemoryProvider {
	return &MemoryProvider{Devices: make(map[uint64]*DeviceProvider)}
}

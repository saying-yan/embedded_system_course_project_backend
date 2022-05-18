package provider

var Provider *MemoryProvider = newMemoryProvider()

// MemoryProvider 为了简化实现，直接将数据存入内存
type MemoryProvider struct {
	Devices map[uint64]*DeviceInfo
}

func (p *MemoryProvider) SetDeviceInfo(d *DeviceInfo) {
	p.Devices[d.ID] = d
}

func (p *MemoryProvider) AddSongs(deviceID uint64, songs []*Song) error {
	device := p.Devices[deviceID]
	for _, s := range songs {
		err := device.AddSong(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func newMemoryProvider() *MemoryProvider {
	return &MemoryProvider{Devices: make(map[uint64]*DeviceInfo)}
}

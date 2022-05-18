package connector

import "errors"

var (
	ErrPacketHeaderInvalid = errors.New("packet header invalid")
	ErrPacketCommand       = errors.New("packet command invalid")
	ErrPacketVersion       = errors.New("packet version invalid")

	ErrPacketDeviceInfo = errors.New("packet device info error")
	ErrPacketSongInfo   = errors.New("packet song info error")
)

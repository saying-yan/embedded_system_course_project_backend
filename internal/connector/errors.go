package connector

import "errors"

var (
	ErrPacketHeaderInvalid = errors.New("packet header invalid")
	ErrPacketCommand       = errors.New("packet command invalid")
	ErrPacketVersion       = errors.New("packet version invalid")
)

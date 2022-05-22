package provider

import "errors"

var (
	ErrSongIDAlreadyExists = errors.New("song ID already exists")
	ErrDeviceNotExists     = errors.New("device not exists")
	ErrSongNotExists       = errors.New("song not exists")
	ErrUnknownListType     = errors.New("unknown list type")
)

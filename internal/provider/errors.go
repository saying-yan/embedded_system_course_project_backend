package provider

import "errors"

var (
	ErrSongIDAlreadyExists = errors.New("song ID already exists")
	ErrDeviceNotExists     = errors.New("device error not exists")
	ErrUnknownListType     = errors.New("unknown list type")
)

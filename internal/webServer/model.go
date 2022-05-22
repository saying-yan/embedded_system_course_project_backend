package web_server

type Song struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	SingerName string `json:"singer_name"`
}

type ListType int

const (
	ListTypeTotal ListType = iota
	ListTypeOrdered
)

type SongListRequest struct {
	ListType ListType `json:"list_type"`
}

type SongListResponse struct {
	Code int `json:"code"`
	Data struct {
		Songs []*Song `json:"songs"`
	} `json:"data"`
}

type OrderSongRequest struct {
	SongID int32 `json:"song_id"`
}

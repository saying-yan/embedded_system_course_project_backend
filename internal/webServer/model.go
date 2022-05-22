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

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewBaseResponse() *BaseResponse {
	return &BaseResponse{}
}

func (resp *BaseResponse) WithCodeOK() *BaseResponse {
	resp.Code = 200
	resp.Msg = "Success"
	return resp
}

func (resp *BaseResponse) WithError(err error) *BaseResponse {
	resp.Code = 500
	resp.Msg = err.Error()
	return resp
}

type SongListRequest struct {
	ListType ListType `json:"list_type"`
}

type SongListResponse struct {
	BaseResponse
	Data struct {
		Songs []*Song `json:"songs"`
	} `json:"data"`
}

type OrderSongRequest struct {
	SongID uint32 `json:"song_id"`
}

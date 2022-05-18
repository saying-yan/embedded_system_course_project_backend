package web_server

type Song struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	SingerName string `json:"singer_name"`
}

type SongListResponse struct {
	Code int `json:"code"`
	Data struct {
		Songs []*Song `json:"songs"`
	} `json:"data"`
}

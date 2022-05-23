package web_server

import (
	"github.com/gin-gonic/gin"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/provider"
	"net/http"
	"strconv"
)

type TestRes struct {
	DeviceID uint64 `json:"device_id"`
}

// TestHandler
// @Summary 测试
// @Schemes
// @Description 测试
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {string} deviceID
// @Router /:deviceID/test [get]
func TestHandler(c *gin.Context) {
	id, _ := c.Get("deviceID")
	deviceID := id.(uint32)
	c.String(http.StatusOK, strconv.Itoa(int(deviceID)))
}

// GetList
// @Summary 获取歌曲列表
// @Schemes
// @Description 获取歌曲列表，0表示全部歌曲的列表，1表示已点歌曲的列表
// @Router /:deviceID/getList [POST]
// @Accept json
// @Param data body SongListRequest true "参数"
// @Success 200 {object} SongListResponse
func GetList(c *gin.Context) {
	id, _ := c.Get("deviceID")
	deviceID := id.(uint32)
	req := SongListRequest{}
	c.BindJSON(&req)

	var songs []*provider.Song
	var err error
	switch req.ListType {
	case ListTypeTotal:
		songs, err = provider.GetDeviceProvider(deviceID).GetList(provider.TotalList)
	case ListTypeOrdered:
		songs, err = provider.GetDeviceProvider(deviceID).GetList(provider.OrderedList)
	default:
		c.JSON(http.StatusOK, NewBaseResponse().WithError(ErrUnknownListType))
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, NewBaseResponse().WithError(err))
		return
	}

	var songsModel []*Song
	for _, song := range songs {
		s := &Song{
			ID:         song.SongID,
			Name:       song.Name,
			SingerName: song.SingerName,
		}
		songsModel = append(songsModel, s)
	}

	resp := SongListResponse{}
	resp.WithCodeOK()
	resp.Data.Songs = songsModel
	c.JSON(http.StatusOK, resp)
}

// OrderSong
// @Summary 点歌
// @Schemes
// @Description 点歌
// @Router /:deviceID/orderSong [POST]
// @Accept json
// @Param data body OrderSongRequest true "参数"
// @Success 200 {object} BaseResponse
func OrderSong(c *gin.Context) {
	id, _ := c.Get("deviceID")
	deviceID := id.(uint32)
	req := OrderSongRequest{}
	c.BindJSON(&req)

	err := provider.GetDeviceProvider(deviceID).OrderSong(req.SongID)
	resp := NewBaseResponse()
	if err != nil {
		resp.WithError(err)
	} else {
		resp.WithCodeOK()
	}
	c.JSON(http.StatusOK, resp)
}

// StickTopSong
// @Summary 置顶
// @Schemes
// @Description 置顶已点歌曲
// @Router /:deviceID/stickTopSong [POST]
// @Accept json
// @Param data body StickTopRequest true "参数"
// @Success 200 {object} BaseResponse
func StickTopSong(c *gin.Context) {
	id, _ := c.Get("deviceID")
	deviceID := id.(uint32)
	req := StickTopRequest{}
	c.BindJSON(&req)

	err := provider.GetDeviceProvider(deviceID).StickTopSong(req.SongIndex)
	resp := NewBaseResponse()
	if err != nil {
		resp.WithError(err)
	} else {
		resp.WithCodeOK()
	}
	c.JSON(http.StatusOK, resp)
}

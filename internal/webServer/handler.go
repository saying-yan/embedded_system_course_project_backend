package web_server

import (
	"github.com/gin-gonic/gin"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/provider"
	"net/http"
)

type TestRes struct {
	DeviceID uint64 `json:"device_id"`
}

func TestHandler(c *gin.Context) {
	id, _ := c.Get("deviceID")
	deviceID := id.(uint64)
	c.JSON(200, TestRes{
		DeviceID: deviceID,
	})

}

func GetList(c *gin.Context) {
	id, _ := c.Get("deviceID")
	deviceID := id.(uint32)
	req := SongListRequest{}
	c.BindJSON(&req)

	var songs []*provider.Song
	var err error
	switch req.ListType {
	case ListTypeTotal:
		songs, err = provider.Provider.GetList(deviceID, provider.TotalList)
	case ListTypeOrdered:
		songs, err = provider.Provider.GetList(deviceID, provider.OrderedList)
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "unknown list type",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
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

	c.JSON(http.StatusOK, SongListResponse{
		Code: 200,
		Data: struct {
			Songs []*Song `json:"songs"`
		}{Songs: songsModel},
	})
}

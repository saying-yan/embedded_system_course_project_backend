package web_server

import (
	"github.com/gin-gonic/gin"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
)

type WebServer struct {
	r    *gin.Engine
	port int
}

func NewWebServer(port int) (*WebServer, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	v1 := r.Group("/v1/:deviceID")
	v1.Use(gin.Recovery(), LoggerMiddleware(), DeviceIDMiddleware())
	{
		v1.GET("/test", TestHandler)
	}

	return &WebServer{
		r:    r,
		port: port,
	}, nil
}

func (web *WebServer) Serve() {
	Logger.Fatal(web.r.Run(":80"))
}

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
	r.Use(LoggerMiddleware(), gin.Recovery())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	return &WebServer{
		r:    r,
		port: port,
	}, nil
}

func (web *WebServer) Serve() {
	Logger.Fatal(web.r.Run(":80"))
}

package web_server

import (
	"github.com/gin-gonic/gin"
	docs "github.com/saying-yan/embedded_system_course_project_backend/docs"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type WebServer struct {
	r    *gin.Engine
	port int
}

func NewWebServer(port int) (*WebServer, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group("/v1/:deviceID")
	v1.Use(gin.Recovery(), LoggerMiddleware(), DeviceIDMiddleware())
	{
		v1.GET("/test", TestHandler)
		v1.POST("/getList", GetList)
		v1.POST("/orderSong", OrderSong)
		v1.POST("/stickTopSong", StickTopSong)
		v1.POST("/nextSong", NextSong)
	}

	return &WebServer{
		r:    r,
		port: port,
	}, nil
}

func (web *WebServer) Serve() {
	Logger.Fatal(web.r.Run(":80"))
}

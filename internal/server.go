package internal

import (
	"fmt"
	config2 "github.com/saying-yan/embedded_system_course_project_backend/internal/config"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	web "github.com/saying-yan/embedded_system_course_project_backend/internal/webServer"
)

type Server struct {
	web *web.WebServer
}

func NewServer(configFile string) (*Server, error) {
	config, err := config2.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("load config %s error: %s\n", configFile, err.Error())
		return nil, err
	}
	err = logger.InitLogger(config.LoggerConf)
	if err != nil {
		fmt.Printf("init logger error: %s", err.Error())
		return nil, err
	}

	w, err := web.NewWebServer(config.WebConf.Port)

	return &Server{
		web: w,
	}, nil
}

func (s *Server) Serve() {
	s.web.Serve()

	select {}
	return
}

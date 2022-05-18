package internal

import (
	"fmt"
	config2 "github.com/saying-yan/embedded_system_course_project_backend/internal/config"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/connector"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	web "github.com/saying-yan/embedded_system_course_project_backend/internal/webServer"
)

type Server struct {
	web       *web.WebServer
	connector *connector.Connector
}

func NewServer(configFile string) (*Server, error) {
	config, err := config2.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("load config %s error: %s\n", configFile, err.Error())
		return nil, err
	}

	err = logger.InitLoggerWithConf(config.LoggerConf)
	if err != nil {
		fmt.Printf("init logger error: %s", err.Error())
		return nil, err
	}

	w, err := web.NewWebServer(config.WebConf.Port)
	if err != nil {
		fmt.Printf("new web server error: %s", err.Error())
		return nil, err
	}

	c, err := connector.NewConnector(config.ConnectorConf.Port)
	if err != nil {
		fmt.Printf("new connector error: %s", err.Error())
		return nil, err
	}

	return &Server{
		web:       w,
		connector: c,
	}, nil
}

func (s *Server) Serve() {
	go func() {
		s.web.Serve()
	}()

	go func() {
		s.connector.Serve()
	}()

	select {}
}

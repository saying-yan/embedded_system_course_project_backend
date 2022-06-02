package connector

import (
	"fmt"
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"net"
	"time"
)

type Connector struct {
	port     int
	ln       net.Listener
	exitChan chan struct{}
}

func NewConnector(port int) (*Connector, error) {
	return &Connector{
		port:     port,
		ln:       nil,
		exitChan: make(chan struct{}),
	}, nil
}

func (c *Connector) Serve() error {
	var err error
	address := fmt.Sprintf(":%d", c.port)
	c.ln, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer c.ln.Close()

	Logger.Infof("connector listen at %s", address)

	go c.checkConnectionActive()

	var tempDelay time.Duration
	for {
		rawConn, err := c.ln.Accept()
		if err != nil {
			Logger.Debugf("accept conn error: %s", err.Error())
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if tempDelay > 1*time.Second {
					tempDelay = 1 * time.Second
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		conn := NewConn(rawConn)
		Logger.Debugf("accept connection from %s", conn.RemoteAddr)
		go conn.handleConn()
	}
}

func (c *Connector) checkConnectionActive() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-c.GetExitChan():
			return
		case <-ticker.C:
			ConnPool.removeTimeoutConn()
		}
	}
}

func (c *Connector) GetExitChan() <-chan struct{} {
	return c.exitChan
}

func (c *Connector) CloseExitChan() {
	ch := c.exitChan
	select {
	case <-ch:
	default:
		close(ch)
	}
}

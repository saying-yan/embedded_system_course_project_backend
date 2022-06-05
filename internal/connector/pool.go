package connector

import (
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"sync"
	"time"
)

var ConnPool = newConnPool()

type ConnectionPool struct {
	connMap         map[uint32]*Conn
	timeoutDuration time.Duration
	rwMutex         sync.RWMutex
}

func newConnPool() *ConnectionPool {
	return &ConnectionPool{
		connMap: make(map[uint32]*Conn),
	}
}

func (pool *ConnectionPool) GetConn(deviceID uint32) *Conn {
	pool.rwMutex.RLock()
	defer pool.rwMutex.RUnlock()

	if conn, ok := pool.connMap[deviceID]; ok {
		return conn
	}
	return nil
}

func (pool *ConnectionPool) PutConn(conn *Conn) {
	pool.rwMutex.Lock()
	defer pool.rwMutex.Unlock()

	pool.connMap[conn.getDeviceID()] = conn
	return
}

func (pool *ConnectionPool) Size() int {
	pool.rwMutex.RLock()
	defer pool.rwMutex.RUnlock()

	return len(pool.connMap)
}

func (pool *ConnectionPool) removeTimeoutConn() {
	pool.rwMutex.Lock()
	defer pool.rwMutex.Unlock()

	if pool.timeoutDuration == 0 {
		// timeoutDuration == 0 means never timeout
		return
	}

	for _, conn := range pool.connMap {
		if time.Now().Sub(conn.getActiveTime()) > pool.timeoutDuration {
			Logger.Debugf("connection:%d from %s timeout", conn.getDeviceID(), conn.RemoteAddr)
			conn.Close()
		}
	}
}

func (pool *ConnectionPool) removeConn(deviceID uint32) {
	pool.rwMutex.Lock()
	defer pool.rwMutex.Unlock()

	if deviceID != 0 {
		delete(pool.connMap, deviceID)
	}
}

package connector

import (
	. "github.com/saying-yan/embedded_system_course_project_backend/internal/logger"
	"sync"
	"time"
)

const (
	TimeoutDuration = 5 * time.Second
)

type ConnPool struct {
	connMap map[uint64]*Conn
	rwMutex sync.RWMutex
}

func newConnPool() *ConnPool {
	return &ConnPool{
		connMap: make(map[uint64]*Conn),
	}
}

func (pool *ConnPool) GetConn(connID uint64) *Conn {
	pool.rwMutex.RLock()
	defer pool.rwMutex.RUnlock()

	if conn, ok := pool.connMap[connID]; ok {
		return conn
	}
	return nil
}

func (pool *ConnPool) PutConn(conn *Conn) {
	pool.rwMutex.Lock()
	defer pool.rwMutex.Unlock()

	pool.connMap[conn.ID] = conn
	return
}

func (pool *ConnPool) Size() int {
	pool.rwMutex.RLock()
	defer pool.rwMutex.RUnlock()

	return len(pool.connMap)
}

func (pool *ConnPool) removeTimeoutConn() {
	pool.rwMutex.Lock()
	defer pool.rwMutex.Unlock()

	for id, conn := range pool.connMap {
		if time.Now().Sub(conn.activeTime) > TimeoutDuration {
			Logger.Debugf("connection:%d from %s timeout", conn.ID, conn.RemoteAddr)
			conn.Close()
			delete(pool.connMap, id)
		}
	}
}

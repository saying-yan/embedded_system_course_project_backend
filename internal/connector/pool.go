package connector

import "sync"

var connPool = newConnPool()

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

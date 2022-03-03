package orbit

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

// Manager 连接管理接口
type Manager interface {
	Add(conn Connection)
	Get(addr string) (Connection, error)
	Len() int
	Del(conn Connection)
	Clear()
}

// manager 连接管理结构体
type manager struct {
	lock  sync.RWMutex
	conns map[string]Connection
}

// Add 添加连接
func (m *manager) Add(conn Connection) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.conns[conn.RemoteAddr()] = conn

	log.Println(fmt.Sprintf("[ MANAGER ] remote addr %s add to connection manager, current connections: %d", conn.RemoteAddr(), m.Len()))
}

// Get 获取当前连接
func (m *manager) Get(addr string) (Connection, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if conn, ok := m.conns[addr]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

// Len 获取当前连接总数
func (m *manager) Len() int {
	return len(m.conns)
}

// Del 关闭并删除连接
func (m *manager) Del(conn Connection) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.conns, conn.RemoteAddr())
	log.Println(fmt.Sprintf("[ MANAGER ] remote addr %s remove from connection manager, current connections: %d", conn.RemoteAddr(), m.Len()))
}

// Clear 清除并停止所有连接
func (m *manager) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for addr, conn := range m.conns {
		conn.Close()
		delete(m.conns, addr)
		log.Println(fmt.Sprintf("[ MANAGER ] remote addr %s remove from connection manager, current connections: %d", conn.RemoteAddr(), m.Len()))
	}

	log.Println(fmt.Sprintf("[ MANAGER ] clear all connections, current connections: %d", m.Len()))
}

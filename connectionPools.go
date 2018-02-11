package main

import (
	"log"
	"net"
	"sync"
)

type connectionPool struct {
	serverName      string
	poolSize        int
	freeConnections []net.Conn
	mux             sync.Mutex
}

func initializePool(poolSize int, serverName string) *connectionPool {
	var connPool connectionPool
	connPool.freeConnections = make([]net.Conn, 0)
	connPool.serverName = serverName
	connPool.poolSize = poolSize
	connPool.addConnections()
	return &connPool
}

func (c *connectionPool) getConnection() net.Conn {
	c.mux.Lock()

	// make new connections if none are free
	if len(c.freeConnections) < 1 {
		c.addConnections()
	}

	conn := c.freeConnections[0]
	c.freeConnections = c.freeConnections[1:]
	c.mux.Unlock()
	return conn
}

func (c *connectionPool) returnConnection(conn net.Conn) {
	c.mux.Lock()
	c.freeConnections = append(c.freeConnections, conn)
	c.mux.Unlock()
}

func (c *connectionPool) addConnections() {
	for i := 0; i < c.poolSize; i++ {
		addr, protocol := configurationServer.GetServerDetails(c.serverName)
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Could not make another connection" + err.Error())
		}
		c.freeConnections = append(c.freeConnections, conn)
	}
}

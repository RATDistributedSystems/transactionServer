package main

import (
	"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	//"strconv"
	//"strings"
	//"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
	"log"
	"sync"
)


type connectionPool struct{
	serverName string
	activeConns int
	connMax int
	freeConnections []net.Conn
	mux sync.Mutex
}


func initializePool(connAmount int, connMax int, serverName string) connectionPool{
	var connPool connectionPool
	connPool.freeConnections = make([]net.Conn, connMax)
	connPool.serverName = serverName
	connPool.connMax = connMax

	i := 0

	//create "connAmount" of connections
	for i < connAmount{
		fmt.Println("loopnumber")
		fmt.Println(i)
		addr, protocol := configurationServer.GetServerDetails(serverName)
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Couldn't Connect to Audit server: " + err.Error())
		}
		connPool.activeConns++
		connPool.freeConnections[i] = conn
		i++
	}
	return connPool
}


func (c * connectionPool) getConnection() net.Conn{
	c.mux.Lock()
	if(len(c.freeConnections) > 0){

		//retrieve first in queue
		conn := c.freeConnections[0]
		//remove first in queue
		c.freeConnections = c.freeConnections[1:]
		//decremement active conns
		c.mux.Unlock()
		return conn

	}else if (len(c.freeConnections) < c.connMax) {
		//if no more usable connections make a new one
		c.openNewConnection()
		//retrieve first in queue
		conn := c.freeConnections[0]
		//remove first in queue
		c.freeConnections = c.freeConnections[1:]
		//decremement active conns
		c.mux.Unlock()
		return conn
	}

	c.mux.Unlock()
	return nil
}

func (c * connectionPool) returnConnection(conn net.Conn){
	c.mux.Lock()
	fmt.Println("returningConnection")
	fmt.Println(len(c.freeConnections))
	c.freeConnections = append(c.freeConnections, conn)
	c.mux.Unlock()
}

func (c * connectionPool) openNewConnection(){
		fmt.Printf("No connections in queue")
		addr, protocol := configurationServer.GetServerDetails(c.serverName)
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Could not make another connection" + err.Error())
		}
		c.activeConns++
		c.freeConnections = append(c.freeConnections, conn)
}
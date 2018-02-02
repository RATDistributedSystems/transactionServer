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
	name string
	activeConns int
	freeConnections []net.Conn
	mux sync.Mutex
	
}

var (
	globalPool connectionPool
	connectionAmount int
	connectionMax int

)

func initializePool(connAmount int, connMax int){
	connectionAmount = connAmount
	connectionMax = connMax
	globalPool.freeConnections = make([]net.Conn, connectionMax)

	i := 0

	//create "connAmount" of connections
	for i < connectionAmount{
		fmt.Println("loopnumber")
		fmt.Println(i)
		addr, protocol := configurationServer.GetServerDetails("audit")
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Couldn't Connect to Audit server: " + err.Error())
		}
		globalPool.activeConns++
		globalPool.freeConnections[i] = conn
		i++
	}
}


func (* connectionPool) getConnection() net.Conn{
	globalPool.mux.Lock()
	if(len(globalPool.freeConnections) > 0){

		//retrieve first in queue
		conn := globalPool.freeConnections[0]
		//remove first in queue
		globalPool.freeConnections = globalPool.freeConnections[1:]
		//decremement active conns
		globalPool.mux.Unlock()
		return conn

	}else if (len(globalPool.freeConnections) < connectionMax) {
		//if no more usable connections make a new one
		openNewConnection()
		//retrieve first in queue
		conn := globalPool.freeConnections[0]
		//remove first in queue
		globalPool.freeConnections = globalPool.freeConnections[1:]
		//decremement active conns
		globalPool.mux.Unlock()
		return conn
	}

	globalPool.mux.Unlock()
	return nil
}

func (* connectionPool) returnConnection(conn net.Conn){
	fmt.Println("returningConnection")
	fmt.Println(len(globalPool.freeConnections))
	globalPool.freeConnections = append(globalPool.freeConnections, conn)
}

func openNewConnection(){
		globalPool.mux.Lock()
		fmt.Printf("No connections in queue")
		addr, protocol := configurationServer.GetServerDetails("audit")
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Could not make another connection" + err.Error())
		}
		globalPool.activeConns++
		globalPool.freeConnections = append(globalPool.freeConnections, conn)
		globalPool.mux.Unlock()
}
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
	//"sync"
)


type connectionPool struct{
	name string
	activeConns int
	freeConnections []net.Conn
	
}

var (
	globalPool connectionPool
	connectionAmount int
	connectionMax int

)

func initializePool(connAmount int, connMax int){
	globalPool.freeConnections = make([]net.Conn, connectionMax)
	connectionAmount = connAmount
	connectionMax = connMax

	i := 0

	//create "connAmount" of connections
	for i < connectionAmount{
		addr, protocol := configurationServer.GetServerDetails("audit")
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Couldn't Connect to Audit server: " + err.Error())
		}
		globalPool.activeConns++
		globalPool.freeConnections[i] = conn
	}
}


func (* connectionPool) getConnection() net.Conn{

	if(len(globalPool.freeConnections) > 0){

		//retrieve first in queue
		conn := globalPool.freeConnections[0]
		//remove first in queue
		globalPool.freeConnections = globalPool.freeConnections[1:]
		//decremement active conns
		return conn

	}else{
		//if no more usable connections make a new one
		openNewConnection()
		//retrieve first in queue
		conn := globalPool.freeConnections[0]
		//remove first in queue
		globalPool.freeConnections = globalPool.freeConnections[1:]
		//decremement active conns
		return conn
	}
}

func (* connectionPool) returnConnection(conn net.Conn){
	globalPool.freeConnections = append(globalPool.freeConnections, conn)
}

func openNewConnection(){
		fmt.Printf("No connections in queue")
		addr, protocol := configurationServer.GetServerDetails("audit")
		conn, err := net.Dial(protocol, addr)
		if err != nil {
			log.Fatalf("Could not make another connection" + err.Error())
		}
		globalPool.activeConns++
		globalPool.freeConnections = append(globalPool.freeConnections, conn)
}
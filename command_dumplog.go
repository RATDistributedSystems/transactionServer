package main

import (
	"net"
	"fmt"
	//"github.com/go-redis/redis"
	//"log"
)

func dumpUser(userId string, filename string){
	fmt.Println("In Dump user")
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "DUMPLOG" + "," + userId + "," + filename
	fmt.Fprintf(conn,text + "\n") 
}

func dump(filename string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "DUMPLOG" + "," + filename
	fmt.Fprintf(conn,text + "\n") 
}
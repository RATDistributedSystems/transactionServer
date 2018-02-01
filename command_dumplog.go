package main

import (
	"net"
	"fmt"
	//"github.com/go-redis/redis"
	//"log"
	"time"
	"strconv"
)

func dumpUser(userId string, filename string, transactionNum int){
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "DUMPLOG", "-1", "", "")
	fmt.Println("In Dump user")
	conn, _ := net.Dial("tcp", "localhost:44445")
	text := "DUMPLOG" + "," + userId + "," + filename
	fmt.Fprintf(conn,text + "\n") 
}

func dump(filename string, transactionNum int){
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "DUMPLOG", "-1", "", "")
	conn, _ := net.Dial("tcp", "localhost:44445")
	text := "DUMPLOG" + "," + filename
	fmt.Fprintf(conn,text + "\n") 
}
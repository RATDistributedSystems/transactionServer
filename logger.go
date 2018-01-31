package main

import "net"
import "fmt"
//import "bufio"
//import "os"
//import "github.com/gocql/gocql"
//import "strconv"
//import "strings"
//import "github.com/go-redis/redis"
//import "github.com/twinj/uuid"
//import "time"
//import "log"



func logUserEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string){
	//connect to audit server
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds
	fmt.Fprintf(conn,text + "\n")
	//go logSystemEvent(time, "AU1", "1",command,userid,"","")
}

func logQuoteEvent(time string, server string, transactionNum string, price string, stockSymbol string, userid string, quoteservertime string, cryptokey string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "Quote" + "," + time + "," + server + "," + transactionNum + "," + price + "," + stockSymbol + "," + userid + "," + quoteservertime + "," + cryptokey
	fmt.Fprintf(conn,text + "\n") 
}

func logSystemEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "System" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds
	fmt.Fprintf(conn,text + "\n") 
}

func logAccountTransactionEvent(time string, server string, transactionNum string, action string, userid string, funds string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "Account" + "," + time + "," + server + "," + transactionNum + "," + action + "," + userid + "," + funds
	fmt.Fprintf(conn,text + "\n") 	
}

func logErrorEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, error string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds + "," + error
	fmt.Fprintf(conn,text + "\n") 
}

func logDebugEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, debug string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds + "," + debug
	fmt.Fprintf(conn,text + "\n") 
}

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


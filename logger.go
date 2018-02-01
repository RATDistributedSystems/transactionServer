package main

//import "net"
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
	//logConnectionect to audit server
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds
	fmt.Fprintf(logConnection,text + "\n")
	//go logSystemEvent(time, "AU1", "1",command,userid,"","")
}

func logQuoteEvent(time string, server string, transactionNum string, price string, stockSymbol string, userid string, quoteservertime string, cryptokey string){
	
	text := "Quote" + "," + time + "," + server + "," + transactionNum + "," + price + "," + stockSymbol + "," + userid + "," + quoteservertime + "," + cryptokey
	fmt.Fprintf(logConnection,text + "\n") 
}

func logSystemEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string){
	
	text := "System" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds
	fmt.Fprintf(logConnection,text + "\n") 
}

func logAccountTransactionEvent(time string, server string, transactionNum string, action string, userid string, funds string){
	
	text := "Account" + "," + time + "," + server + "," + transactionNum + "," + action + "," + userid + "," + funds
	fmt.Fprintf(logConnection,text + "\n") 	
}

func logErrorEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, error string){
	
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds + "," + error
	fmt.Fprintf(logConnection,text + "\n") 
}

func logDebugEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, debug string){
	
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds + "," + debug
	fmt.Fprintf(logConnection,text + "\n") 
}



package main

import (
	"bufio"
	"fmt"
	"strings"
)

func quoteCacheRequest(userId string, stockSymbol string, transactionNum int) []string {
	//log events
	logUserEvent("TS1", transactionNum, "QUOTE", userId, stockSymbol, "")

	//check if a cached quote exists
	if(cacheExists(stockSymbol) == true){
		//obtain values 
		return cacheReturn(stockSymbol)
	}else{
		//if it doesnt access the quote server normally
		conn := GetQuoteServerConnection() //conn := quotePool.getConnection()
		fmt.Fprintf(conn, "%s,%s", stockSymbol, userId)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		conn.Close() //quotePool.returnConnection(conn)
		messageArray := strings.Split(message, ",")
		logQuoteEvent("TS1", transactionNum, messageArray[0], messageArray[1], messageArray[2], messageArray[3], messageArray[4])
		return messageArray
	}
}

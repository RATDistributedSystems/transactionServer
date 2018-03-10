package main

import (
	"bufio"
	"fmt"
	"strings"
)

func quoteRequest(userID string, stockSymbol string, transactionNum int) int {
	//logUserEvent("TS1", transactionNum, "QUOTE", userID, stockSymbol, "")

	// Make Quote Request
	conn := GetQuoteServerConnection()
	fmt.Fprintf(conn, "%s,%s\n", stockSymbol, userID)

	// Get response
	message, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	messageArray := strings.Split(message, ",")
	logQuoteEvent(serverName, transactionNum, messageArray[0], messageArray[1], messageArray[2], messageArray[3], messageArray[4])
	return stringToCents(messageArray[0])
}

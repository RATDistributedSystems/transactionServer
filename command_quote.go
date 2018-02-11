package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func quoteRequest(userId string, stockSymbol string, transactionNum int) []string {
	conn := quotePool.getConnection()
	stockSymbol = strings.TrimSuffix(stockSymbol, "\n")
	userId = strings.TrimSuffix(userId, "\n")
	text := stockSymbol + "," + userId

	fmt.Fprintf(conn, text+"\n")
	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	//fmt.Print("Message from server: " + message)
	quotePool.returnConnection(conn)
	messageArray := strings.Split(message, ",")

	timestamp_q := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_quote := strconv.FormatInt(timestamp_q, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum), 10)
	logQuoteEvent(timestamp_quote, "TS1", transactionNum_string, messageArray[0], messageArray[1], messageArray[2], messageArray[3], messageArray[4])

	return messageArray
}

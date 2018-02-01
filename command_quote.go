package main

import (
	"net"
	"fmt"
	"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	"strings"
	//"github.com/twinj/uuid"
	"time"
	//"github.com/go-redis/redis"
	//"log"
)

func quoteRequest(userId string, stockSymbol string ,transactionNum int) []string{
	// connect to this socket
		conn, _ := net.Dial("tcp", "localhost:3333")
		stockSymbol = strings.TrimSuffix(stockSymbol, "\n")
		userId = strings.TrimSuffix(userId, "\n")
		text := stockSymbol + "," + userId
		fmt.Print(text)
		// send to socket
		fmt.Fprintf(conn, text + "\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)

		messageArray := strings.Split(message, ",")

		timestamp_q := (time.Now().UTC().UnixNano())/ 1000000
		timestamp_quote := strconv.FormatInt(timestamp_q,10)
		transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
		logQuoteEvent(timestamp_quote,"TS1",transactionNum_string,messageArray[0],messageArray[1],userId,messageArray[3],messageArray[4])
		timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
		timestamp_command := strconv.FormatInt(timestamp_time, 10)
		//transactionNum_user += 1
		//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
		logUserEvent(timestamp_command, "TS1", transactionNum_string, "QUOTE", userId, messageArray[1], "")

		return messageArray
}
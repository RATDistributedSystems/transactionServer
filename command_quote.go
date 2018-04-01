package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type commandQuote struct {
	username string
	stock    string
}

func (c commandQuote) process(transaction int) string {
	logUserEvent(serverName, transaction, "QUOTE", c.username, c.stock, "")
	return centsToString(getQuote(c.username, c.stock, transaction))
}

func getQuoteServerConnection() net.Conn {
	addr, protocol := configurationServer.GetServerDetails("quote")
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Printf("Encountered error when trying to connect to quote server\n%s", err.Error())
	}
	return conn
}

func getQuote(userID string, stockSymbol string, transactionNum int) int {
	exists, quote := cacheExists(stockSymbol)
	//check if a cached quote exists
	if exists {
		return quote
	}

	conn := getQuoteServerConnection()
	fmt.Fprintf(conn, "%s,%s \n", stockSymbol, userID)
	msg, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	cacheAdd(stockSymbol, msg)
	msgSplit := strings.Split(msg, ",")
	logQuoteEvent(serverName, transactionNum, msgSplit[0], msgSplit[1], msgSplit[2], msgSplit[3], msgSplit[4])
	return stringToCents(msgSplit[0])
}

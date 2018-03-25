package main

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/RATDistributedSystems/utilities"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
	"github.com/gocql/gocql"
	"github.com/twinj/uuid"
)

var sessionGlobal *gocql.Session
var transactionNumGlobal = 0
var auditPool *connectionPool
var configurationServer = utilities.Load()
var serverName string

//var quotePool = initializePool(100, "quote")

func main() {
	configurationServer.Pause()
	uuid.Init()
	serverName = uuid.NewV4().String()
	initCassandra()
	auditPool = initializePool(150, 190, "audit")
	initRedis()
	initTCPListener()

}

func GetQuoteServerConnection() net.Conn {
	addr, protocol := configurationServer.GetServerDetails("quote")
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Printf("Encountered error when trying to connect to quote server\n%s", err.Error())
	}
	return conn
}

func initCassandra() {
	//connect to database
	hostname := configurationServer.GetValue("transdb_ip")
	keyspace := configurationServer.GetValue("transdb_keyspace")
	protocol := configurationServer.GetValue("transdb_proto")
	ratdatabase.InitCassandraConnection(hostname, keyspace, protocol)
	sessionGlobal = ratdatabase.CassandraConnection
}

func initTCPListener() {
	// Listen for incoming connections.

	addr, protocol := configurationServer.GetListnerDetails("transaction")
	l, err := net.Listen(protocol, addr)
	if err != nil {
		panic(err)
	}
	//defer l.Close()
	log.Printf("Transaction server listening on %s", addr)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Panic(err)
		} else {
			transactionNumGlobal++
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, transactionNumGlobal)
	}

}

func handleRequest(conn net.Conn, transactionNumGlobal int) {
	tp := textproto.NewReader(bufio.NewReader(conn))
	untrimmedMsg, err := tp.ReadLine()
	if err != nil {
		log.Println(err)
	}
	conn.Close()
	message := strings.TrimSpace(untrimmedMsg)
	executeCommand(message, transactionNumGlobal)
}

func executeCommand(command string, transactionNumGlobal int) {
	result := strings.Split(command, ",")
	log.Printf("Recieved Request: [%d]  %s", transactionNumGlobal, command)

	//strconv.FormatInt(int64(buyableStocks), 10)
	//strconv.Atoi

	switch result[0] {
	case "ADD":
		var x, _ = strconv.Atoi(result[3])
		logUserEvent(serverName, x, "ADD", result[1], "", result[2])
		addUser(result[1], result[2], x)
	case "QUOTE":
		var x, _ = strconv.Atoi(result[3])
		logUserEvent(serverName, x, "QUOTE", result[1], result[2], "")
		//quoteRequest(result[1], result[2], transactionNumGlobal)
		quoteCacheRequest(result[1], result[2], x)
	case "BUY":
		var x, _ = strconv.Atoi(result[4])
		logUserEvent(serverName, x, "BUY", result[1], result[2], result[3])
		buy(result[1], result[2], result[3], x)
	case "COMMIT_BUY":
		var x, _ = strconv.Atoi(result[2])
		logUserEvent(serverName, x, "COMMIT_BUY", result[1], "", "")
		commitBuy(result[1], x)
	case "CANCEL_BUY":
		var x, _ = strconv.Atoi(result[2])
		logUserEvent(serverName, x, "CANCEL_BUY", result[1], "", "")
		cancelBuy(result[1], x)
	case "SELL":
		var x, _ = strconv.Atoi(result[4])
		logUserEvent(serverName, x, "SELL", result[1], result[2], result[3])
		sell(result[1], result[2], result[3], x)
	case "COMMIT_SELL":
		var x, _ = strconv.Atoi(result[2])
		logUserEvent(serverName, x, "COMMIT_SELL", result[1], "", "")
		commitSell(result[1], x)
	case "CANCEL_SELL":
		var x, _ = strconv.Atoi(result[2])
		logUserEvent(serverName, x, "CANCEL_SELL", result[1], "", "")
		cancelSell(result[1], x)
	case "SET_BUY_AMOUNT":
		var x, _ = strconv.Atoi(result[4])
		logUserEvent(serverName, x, "SET_BUY_AMOUNT", result[1], result[2], result[3])
		//setBuyAmount(result[1], result[2], result[3], x)
	case "SET_BUY_TRIGGER":
		var x, _ = strconv.Atoi(result[4])
		logUserEvent(serverName, x, "SET_BUY_TRIGGER", result[1], result[2], result[3])
		//setBuyTrigger(result[1], result[2], result[3], x)
	case "CANCEL_SET_BUY":
		var x, _ = strconv.Atoi(result[3])
		logUserEvent(serverName, x, "CANCEL_SET_BUY", result[1], result[2], "")
		//cancelBuyTrigger(result[1], result[2], x)
	case "SET_SELL_AMOUNT":
		var x, _ = strconv.Atoi(result[4])
		logUserEvent(serverName, x, "SET_SELL_AMOUNT", result[1], result[2], result[3])
		//setSellAmount(result[1], result[2], result[3], x)
	case "SET_SELL_TRIGGER":
		var x, _ = strconv.Atoi(result[4])
		logUserEvent(serverName, x, "SET_SELL_TRIGGER", result[1], result[2], result[3])
		//setSellTrigger(result[1], result[2], result[3], x)
	case "CANCEL_SET_SELL":
		var x, _ = strconv.Atoi(result[3])
		logUserEvent(serverName, x, "CANCEL_SET_SELL", result[1], result[2], "")
		//cancelSellTrigger(result[1], result[2], x)
	case "DISPLAY_SUMMARY":
		var x, _ = strconv.Atoi(result[2])
		logUserEvent(serverName, x, "DISPLAY_SUMMARY", result[1], "", "")
		//displaySummary(result[1], x)
	case "DUMPLOG":

		if len(result) == 4 {
			var x, _ = strconv.Atoi(result[3])
			logUserEvent(serverName, x, "DUMPLOG", result[1], "", "")
			dumpUser(result[1], result[2], x)
		} else if len(result) == 3 {
			var x, _ = strconv.Atoi(result[2])
			logUserEvent(serverName, x, "DUMPLOG", "-1", "", "")
			dump(result[1], x)
		}
	}

}

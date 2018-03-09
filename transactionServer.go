package main

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"strings"

	"github.com/RATDistributedSystems/utilities"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
	"github.com/gocql/gocql"
	"github.com/twinj/uuid"
)

var sessionGlobal *gocql.Session
var transactionNumGlobal = 0
var configurationServer = utilities.GetConfigurationFile("config.json")
var auditPool = initializePool(150, 190, "audit")
var serverName = configurationServer.GetValue("ts_name")

//var quotePool = initializePool(100, "quote")

func main() {
	uuid.Init()
	initCassandra()
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
	hostname := configurationServer.GetValue("cassandra_ip")
	keyspace := configurationServer.GetValue("cassandra_keyspace")
	protocol := configurationServer.GetValue("cassandra_proto")
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
		} else{
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

	switch result[0] {
	case "ADD":
		logUserEvent(serverName, transactionNumGlobal, "ADD", result[1], "", result[2])
		addUser(result[1], result[2], result[3])
	case "QUOTE":
		logUserEvent(serverName, transactionNumGlobal, "QUOTE", result[1], result[2], "")
		//quoteRequest(result[1], result[2], transactionNumGlobal)
		quoteCacheRequest(result[1], result[2], result[3])
	case "BUY":
		logUserEvent(serverName, transactionNumGlobal, "BUY", result[1], result[2], result[3])
		buy(result[1], result[2], result[3], result[4])
	case "COMMIT_BUY":
		logUserEvent(serverName, result[2], "COMMIT_BUY", result[1], "", "")
		commitBuy(result[1], result[2])
	case "CANCEL_BUY":
		logUserEvent(serverName, result[2], "CANCEL_BUY", result[1], "", "")
		cancelBuy(result[1], result[2])
	case "SELL":
		logUserEvent(serverName, result[4], "SELL", result[1], result[2], result[3])
		sell(result[1], result[2], result[3], result[])
	case "COMMIT_SELL":
		logUserEvent(serverName, result[2], "COMMIT_SELL", result[1],"", "")
		commitSell(result[1], result[2])
	case "CANCEL_SELL":
		logUserEvent(serverName, result[2], "CANCEL_SELL", result[1], "", "")
		cancelSell(result[1], result[2])
	case "SET_BUY_AMOUNT":
		logUserEvent(serverName, result[4], "SET_BUY_AMOUNT", result[1], result[2], result[3])
		setBuyAmount(result[1], result[2], result[3], result[4])
	case "SET_BUY_TRIGGER":
		logUserEvent(serverName, result[4], "SET_BUY_TRIGGER", result[1], result[2], result[3])
		setBuyTrigger(result[1], result[2], result[3], result[4])
	case "CANCEL_SET_BUY":
		logUserEvent(serverName, result[3], "CANCEL_SET_BUY", result[1], result[2], "")
		cancelBuyTrigger(result[1], result[2], result[3])
	case "SET_SELL_AMOUNT":
		logUserEvent(serverName, result[4], "SET_SELL_AMOUNT", result[1], result[2], result[3])
		setSellAmount(result[1], result[2], result[3], result[4])
	case "SET_SELL_TRIGGER":
		logUserEvent(serverName, result[4], "SET_SELL_TRIGGER", result[1], result[2], result[3])
		setSellTrigger(result[1], result[2], result[3], result[4])
	case "CANCEL_SET_SELL":
		logUserEvent(serverName, result[3], "CANCEL_SET_SELL", result[1], result[2], "")
		cancelSellTrigger(result[1], result[2], result[3])
	case "DISPLAY_SUMMARY":
		logUserEvent(serverName, result[2], "DISPLAY_SUMMARY", result[1], "", "")
		displaySummary(result[1], result[2])
	case "DUMPLOG":
		if len(result) == 3 {
			logUserEvent(serverName, result[3], "DUMPLOG", result[1], "", "")
			dumpUser(result[1], result[2], result[3])
		} else if len(result) == 2 {
			logUserEvent(serverName, result[2], "DUMPLOG", "-1", "", "")
			dump(result[1], result[2])
		}
	}

}

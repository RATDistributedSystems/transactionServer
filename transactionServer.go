package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/RATDistributedSystems/utilities"
	"github.com/gocql/gocql"
	"github.com/twinj/uuid"
)

var sessionGlobal *gocql.Session

var transactionNumGlobal = 0
var configurationServer = utilities.GetConfigurationFile("config.json")
var auditPool = initializePool(100, "audit")

//var quotePool = initializePool(100, "quote")

func main() {
	uuid.Init()
	initCassandra()
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
	cluster := gocql.NewCluster(configurationServer.GetValue("cassandra_ip"))
	cluster.Keyspace = configurationServer.GetValue("cassandra_keyspace")
	proto, err := strconv.Atoi(configurationServer.GetValue("cassandra_proto"))
	if err != nil {
		panic("Cassandra protocol version not int")
	}
	cluster.ProtoVersion = proto
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	sessionGlobal = session
	fmt.Println("Database Connection Created")
}

func initTCPListener() {
	// Listen for incoming connections.
	addr, protocol := configurationServer.GetServerDetails("transaction")
	l, err := net.Listen(protocol, addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	log.Printf("Transaction server listening on %s", addr)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	tp := textproto.NewReader(bufio.NewReader(conn))
	untrimmedMsg, err := tp.ReadLine()
	if err != nil {
		log.Println(err)
	}
	message := strings.TrimSpace(untrimmedMsg)
	log.Printf("Recieved Request: [%d]  %s", transactionNumGlobal, message)
	commandExecuter(message)
}

func commandExecuter(command string) {
	result := strings.Split(command, ",")
	transactionNumGlobal++

	switch result[0] {
	case "ADD":
		logUserEvent("TS1", transactionNumGlobal, "ADD", result[1], "", result[2])
		addUser(result[1], result[2], transactionNumGlobal)
	case "QUOTE":
		logUserEvent("TS1", transactionNumGlobal, "QUOTE", result[1], result[2], "")
		quoteRequest(result[1], result[2], transactionNumGlobal)
	case "BUY":
		logUserEvent("TS1", transactionNumGlobal, "BUY", result[1], result[2], result[3])
		buy(result[1], result[2], result[3], transactionNumGlobal)
	case "COMMIT_BUY":
		logUserEvent("TS1", transactionNumGlobal, "COMMIT_BUY", result[1], "", "")
		commitBuy(result[1], transactionNumGlobal)
	case "CANCEL_BUY":
		logUserEvent("TS1", transactionNumGlobal, "CANCEL_BUY", result[1], "", "")
		cancelBuy(result[1], transactionNumGlobal)
	case "SELL":
		logUserEvent("TS1", transactionNumGlobal, "SELL", result[1], result[2], result[3])
		sell(result[1], result[2], result[3], transactionNumGlobal)
	case "COMMIT_SELL":
		logUserEvent("TS1", transactionNumGlobal, "COMMIT_SELL", result[1], "", "")
		commitSell(result[1], transactionNumGlobal)
	case "CANCEL_SELL":
		logUserEvent("TS1", transactionNumGlobal, "CANCEL_SELL", result[1], "", "")
		cancelSell(result[1], transactionNumGlobal)
	case "SET_BUY_AMOUNT":
		logUserEvent("TS1", transactionNumGlobal, "SET_BUY_AMOUNT", result[1], result[2], result[3])
		setBuyAmount(result[1], result[2], result[3], transactionNumGlobal)
	case "SET_BUY_TRIGGER":
		logUserEvent("TS1", transactionNumGlobal, "SET_BUY_TRIGGER", result[1], result[2], result[3])
		setBuyTrigger(result[1], result[2], result[3], transactionNumGlobal)
	case "CANCEL_SET_BUY":
		logUserEvent("TS1", transactionNumGlobal, "CANCEL_SET_BUY", result[1], result[2], "")
		cancelBuyTrigger(result[1], result[2], transactionNumGlobal)
	case "SET_SELL_AMOUNT":
		logUserEvent("TS1", transactionNumGlobal, "SET_SELL_AMOUNT", result[1], result[2], result[3])
		setSellAmount(result[1], result[2], result[3], transactionNumGlobal)
	case "SET_SELL_TRIGGER":
		logUserEvent("TS1", transactionNumGlobal, "SET_SELL_TRIGGER", result[1], result[2], result[3])
		setSellTrigger(result[1], result[2], result[3], transactionNumGlobal)
	case "CANCEL_SET_SELL":
		logUserEvent("TS1", transactionNumGlobal, "CANCEL_SET_SELL", result[1], result[2], "")
		cancelSellTrigger(result[1], result[2], transactionNumGlobal)
	case "DISPLAY_SUMMARY":
		logUserEvent("TS1", transactionNumGlobal, "DISPLAY_SUMMARY", result[1], "", "")
		displaySummary(result[1], transactionNumGlobal)
	case "DUMPLOG":
		if len(result) == 3 {
			logUserEvent("TS1", transactionNumGlobal, "DUMPLOG", result[1], "", "")
			dumpUser(result[1], result[2], transactionNumGlobal)
		} else if len(result) == 2 {
			logUserEvent("TS1", transactionNumGlobal, "DUMPLOG", "-1", "", "")
			dump(result[1], transactionNumGlobal)
		}
	}

}

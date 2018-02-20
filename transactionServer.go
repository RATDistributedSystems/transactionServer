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
var auditPool = initializePool(100, 120, "audit")

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
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	tp := textproto.NewReader(bufio.NewReader(conn))
	untrimmedMsg, err := tp.ReadLine()
	if err != nil {
		log.Println(err)
	}
	conn.Close()
	message := strings.TrimSpace(untrimmedMsg)
	executeCommand(message)
}

func executeCommand(command string) {
	result := strings.Split(command, ",")
	transactionNumGlobal++
	log.Printf("Recieved Request: [%d]  %s", transactionNumGlobal, command)

	switch result[0] {
	case "ADD":
		addUser(result[1], result[2], transactionNumGlobal)
	case "QUOTE":
		quoteRequest(result[1], result[2], transactionNumGlobal)
	case "BUY":
		buy(result[1], result[2], result[3], transactionNumGlobal)
	case "COMMIT_BUY":
		commitBuy(result[1], transactionNumGlobal)
	case "CANCEL_BUY":
		cancelBuy(result[1], transactionNumGlobal)
	case "SELL":
		sell(result[1], result[2], result[3], transactionNumGlobal)
	case "COMMIT_SELL":
		commitSell(result[1], transactionNumGlobal)
	case "CANCEL_SELL":
		cancelSell(result[1], transactionNumGlobal)
	case "SET_BUY_AMOUNT":
		setBuyAmount(result[1], result[2], result[3], transactionNumGlobal)
	case "SET_BUY_TRIGGER":
		setBuyTrigger(result[1], result[2], result[3], transactionNumGlobal)
	case "CANCEL_SET_BUY":
		cancelBuyTrigger(result[1], result[2], transactionNumGlobal)
	case "SET_SELL_AMOUNT":
		setSellAmount(result[1], result[2], result[3], transactionNumGlobal)
	case "SET_SELL_TRIGGER":
		setSellTrigger(result[1], result[2], result[3], transactionNumGlobal)
	case "CANCEL_SET_SELL":
		cancelSellTrigger(result[1], result[2], transactionNumGlobal)
	case "DISPLAY_SUMMARY":
		displaySummary(result[1], transactionNumGlobal)
	case "DUMPLOG":
		if len(result) == 3 {
			dumpUser(result[1], result[2], transactionNumGlobal)
		} else if len(result) == 2 {
			dump(result[1], transactionNumGlobal)
		}
	}

}

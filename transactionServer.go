package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/RATDistributedSystems/utilities"
	"github.com/gocql/gocql"
	"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
)

var sessionGlobal *gocql.Session
var logConnection net.Conn
var transactionNumGlobal int
var configurationServer = utilities.GetConfigurationFile("config.json")

func main() {
	transactionNumGlobal = 0
	initServer()
	initAuditConnection()
	uuid.Init()
	tcpListener()
}

func initAuditConnection(){
	conn, _ := net.Dial("tcp", "localhost:44445")
	logConnection = conn
}

func initServer() {
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

func tcpListener() {
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
	message, _ := bufio.NewReader(conn).ReadString('\n')
	//remove new line character and any spaces received
	message = strings.TrimSuffix(message, "\n")
	message = strings.TrimSpace(message)
	log.Printf("Recieved Request: %s", message)
	commandExecuter(message)
	conn.Close()
}

func commandExecuter(command string) {
	result := strings.Split(command, ",")
	//incrementing here since workload has no invalid entries
	transactionNumGlobal++

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

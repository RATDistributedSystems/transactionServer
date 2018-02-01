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
var transactionNumGlobal int
var configurationServer = utilities.GetConfigurationFile("config.json")

func main() {
	initServer()
	uuid.Init()
	tcpListener()
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
	transactionNumGlobal = 0
	fmt.Println("Database Connection Created")
}

func tcpListener() {
	// Listen for incoming connections.
	addr, protocol := configurationServer.GetServerDetails("listener")
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

func stringToCents(x string) int {
	result := strings.Split(x, ".")
	dollars, err := strconv.Atoi(result[0])
	if err != nil {
		log.Printf("Couldn't convert %s to int", result[0])
		return 0
	}

	cents, err := strconv.Atoi(strings.TrimSuffix(result[1], "\n"))
	if err != nil {
		log.Printf("Couldn't convert %s to int", result[1])
		return 0
	}

	return (dollars * 100) + cents
}

//chekcs if the command can be executed
//ie check if buy set before commit etc
func checkDependency(command string, userId string, stock string) bool {

	count := 0
	var err error = nil
	switch command {
	case "COMMIT_BUY":
		err = sessionGlobal.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "COMMIT_SELL":
		err = sessionGlobal.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "CANCEL_BUY":
		sessionGlobal.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "CANCEL_SELL":
		err = sessionGlobal.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "CANCEL_SET_BUY":
		err = sessionGlobal.Query("SELECT count(*) FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count)
	case "CANCEL_SET_SELL":
		err = sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count)
	}

	if err != nil {
		panic(err)
	}
	return count != 0
}

func getUsableCash(userId string) int {
	var usableCash int
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(err)
	}

	return usableCash
}

func addFunds(userId string, addCashAmount int) {
	usableCash := getUsableCash(userId)
	totalCash := usableCash + addCashAmount
	totalCashString := strconv.FormatInt(int64(totalCash), 10)

	//return add funds to user
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(err)
	}

}

//check if the trigger hasn't been cancelled
func checkTriggerExists(userId string, stock string, isBuyOperation bool) bool {

	var count int

	if isBuyOperation == true {
		if err := sessionGlobal.Query("SELECT count(*) FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil {
			panic(err)
		}
	} else {
		if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil {
			panic(err)
		}
	}

	return count == 1
}

func checkStockOwnership(userId string, stock string) (int, string) {

	var ownedstockname string
	var ownedstockamount int
	var usid string
	//var hasStock bool

	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
	for iter.Scan(&usid, &ownedstockname, &ownedstockamount) {
		if ownedstockname == stock {
			//hasStock = true
			break
		}
	}
	if err := iter.Close(); err != nil {
		panic(err)
	}

	//returns 0 if stock is empty
	return ownedstockamount, usid

}

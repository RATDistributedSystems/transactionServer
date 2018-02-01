package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"github.com/gocql/gocql"
	"strings"
	"github.com/twinj/uuid"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3334"
    CONN_TYPE = "tcp"

    DB_HOST = "localhost"
    DB_KEYSPACE = "userdb"
)

var sessionGlobal *gocql.Session
var transactionNumGlobal int

func main(){
	transactionNumGlobal = 0;
	initializeDbConnection()
	uuid.Init()
	startListening();
}

func initializeDbConnection(){
	//connect to database
	cluster := gocql.NewCluster(DB_HOST)
	cluster.Keyspace = DB_KEYSPACE
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	sessionGlobal = session
	fmt.Println("Database Connection Created")
}


func startListening(){
	// Listen for incoming connections.
	fmt.Println("TCP listner started")
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		//<---------------------------NEW THREAD PER REQUEST----------------------------->>
		go handleRequest(conn)
	}
	// Close the listener when the application closes.
	defer l.Close()
}

func handleRequest(conn net.Conn) {
  // will listen for message to process ending in newline (\n)
    print("received request")
    message, _ := bufio.NewReader(conn).ReadString('\n')
    //remove new line character and any spaces received
    message = strings.TrimSuffix(message, "\n")
    message = strings.TrimSpace(message)
	commandExecuter(message)
  	conn.Close()
}


func commandExecuter(command string){
		result := processCommand(command)
		caseValue := result[0]
		//incrementing here since workload has no invalid entries
		transactionNumGlobal++

		switch caseValue {
			case "ADD":
				addUser(result[1],result[2],transactionNumGlobal)

			case "QUOTE":
				quoteRequest(result[1],result[2],transactionNumGlobal)

			case "BUY":
				buy(result[1],result[2],result[3],transactionNumGlobal)

			case "COMMIT_BUY":
				commitBuy(result[1],transactionNumGlobal)

			case "CANCEL_BUY":
				cancelBuy(result[1],transactionNumGlobal)

			case "SELL":
				sell(result[1],result[2],result[3],transactionNumGlobal)

			case "COMMIT_SELL":
				commitSell(result[1],transactionNumGlobal)

			case "CANCEL_SELL":
				cancelSell(result[1],transactionNumGlobal)

			case "SET_BUY_AMOUNT":
				setBuyAmount(result[1],result[2],result[3],transactionNumGlobal)

			case "SET_BUY_TRIGGER":
				setBuyTrigger(result[1],result[2],result[3],transactionNumGlobal)

			case "CANCEL_SET_BUY":
				cancelBuyTrigger(result[1],result[2],transactionNumGlobal)
		
			case "SET_SELL_AMOUNT":
				setSellAmount(result[1],result[2],result[3],transactionNumGlobal)
		
			case "SET_SELL_TRIGGER":
				setSellTrigger(result[1],result[2],result[3],transactionNumGlobal)
		
			case "CANCEL_SET_SELL":
				cancelSellTrigger(result[1],result[2],transactionNumGlobal)
		
			case "DISPLAY_SUMMARY":
				displaySummary(result[1], transactionNumGlobal)

			case "DUMPLOG":
				if len(result) == 3{
					dumpUser(result[1],result[2])
				} else if len(result) == 2{
					dump(result[1])
				}
		}
}




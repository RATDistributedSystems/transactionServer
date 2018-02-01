package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"github.com/gocql/gocql"
	"strconv"
	"strings"
	"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
	//"log"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3334"
    CONN_TYPE = "tcp"
)

var sessionGlobal *gocql.Session
var transactionNumGlobal int

func main(){
	initializeServer()
	uuid.Init()
	tcpListener();
	//commandListener();
}

func initializeServer(){
	//connect to database
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	sessionGlobal = session
	transactionNumGlobal = 0;
	fmt.Println("Database Connection Created")
}


func tcpListener(){
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

func getUsableCash(userId string) int{
	var usableCash int
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	return usableCash
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
/*
func commandListener(){
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a command: ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)

		result := processCommand(text)

		if result[0] == "ADD"{
			addUser(result[1],result[2])
		}

		if result[0] == "QUOTE"{
			quoteRequest(result[1],result[2])
		}

		if result[0] == "BUY"{
			buy(result[1],result[2],result[3])
		}

		if result[0] == "BUY_COMMIT"{
			commitBuy(result[1])
		}

		if result[0] == "CANCEL_BUY"{
			cancelBuy(result[1])
		}

		if result[0] == "SELL"{
			sell(result[1],result[2],result[3])
		}

		if result[0] == "SELL_COMMIT"{
			fmt.Println(len(result))
			commitSell(result[1])
		}

		if result[0] == "CANCEL_SELL"{
			cancelSell(result[1])
		}

		if result[0] == "SET_BUY_AMOUNT"{
			fmt.Println(len(result))
			setBuyAmount(result[1],result[2],result[3])
		}

		if result[0] == "SET_BUY_TRIGGER"{
			fmt.Println(len(result))
			setBuyTrigger(result[1],result[2],result[3])
		}

		if result[0] == "SET_SELL_AMOUNT"{
			fmt.Println(len(result))
			setSellAmount(result[1],result[2],result[3])
		}

		if result[0] == "SET_SELL_TRIGGER"{
			fmt.Println(len(result))
			setSellTrigger(result[1],result[2],result[3])
		}
		if result[0] == "CANCEL_SELL_TRIGGER"{
			fmt.Println(len(result))
			cancelSellTrigger(result[1],result[2])
		}
		if result[0] == "CANCEL_BUY_TRIGGER"{
			fmt.Println(len(result))
			cancelBuyTrigger(result[1],result[2])
		}

		if result[0] == "DUMPLOG"{
			if len(result) == 3{
				dumpUser(result[1],result[2])
			} else if len(result) == 2{
				dump(result[1])
			}
		}


	}
}
*/
func processCommand(text string) []string{
	result := strings.Split(text, ",")
	for i := range result {
		fmt.Println(result[i])
	}
	return result;
}




func stringToCents(x string) int{

	var dollars int
	var cents int

	fmt.Println(x)
	result := strings.Split(x, ".")
	for i := range result {
		fmt.Println(result[i])
	}

	if i, err := strconv.Atoi(result[0]); err == nil {
		dollars = i
		fmt.Println("dollar converted to int")
		fmt.Println(i)
	}

	result[1] = strings.TrimSuffix(result[1], "\n")
	if e, err := strconv.Atoi(result[1]); err == nil {
		cents = e
		fmt.Println("cents converted to int")
		fmt.Println(e)
	}

	dollars = dollars * 100
	var money int = dollars + cents

	return money
}

//chekcs if the command can be executed
//ie check if buy set before commit etc
func checkDependency(command string, userId string, stock string) bool{

	var count int

	if command == "COMMIT_BUY"{
		//check if a buy entry exists in buypendingtransactions table

		if err := sessionGlobal.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
		if count == 0{
			return false

		}else{

			return true
		}
	}
	if command == "COMMIT_SELL"{
		//check if a sell entry exists in sellpendingtransactions table
			//check if a sell entry exists in buypendingtransactions table

			if err := sessionGlobal.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}
			if count == 0{
				return false
			}else{
				return true
			}
	}
	if command == "CANCEL_BUY"{
		//check if a buy entry exists in buypendingtransactions table
		if err := sessionGlobal.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}

		
	}
	if command == "CANCEL_SELL"{
		//check if a sell entry exists in sellpendingtransactions table
		if err := sessionGlobal.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}
		
	}
	if command == "CANCEL_SET_BUY"{
		if err := sessionGlobal.Query("SELECT count(*) FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}
	}
	if command == "CANCEL_SET_SELL"{
		if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}
		
	}

	return false
}

func addFunds(userId string, addCashAmount int){

	var usableCash int

	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	totalCash := usableCash + addCashAmount;
	totalCashString := strconv.FormatInt(int64(totalCash), 10)

	//return add funds to user
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}








//check if the trigger hasn't been cancelled
func checkTriggerExists(userId string, stock string, operation bool) bool{


	var count int

	if operation == true {
		if err := sessionGlobal.Query("SELECT count(*) FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
	}else{
		if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
	}

	if count == 1 {
		return true
	}else{
		return false
	}
}









func checkStockOwnership(userId string, stock string) (int, string){

	var ownedstockname string
	var ownedstockamount int
	var usid string
	//var hasStock bool

	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='"+ userId + "'").Iter()
	for iter.Scan(&usid, &ownedstockname, &ownedstockamount) {
		if (ownedstockname == stock){
			//hasStock = true
			break;
		}
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//returns 0 if stock is empty
	return ownedstockamount, usid
	
}


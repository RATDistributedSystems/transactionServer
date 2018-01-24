package main

import "net"
import "fmt"
import "bufio"
import "os"
import "github.com/gocql/gocql"
import "strconv"
import "strings"
//import "github.com/go-redis/redis"
import "github.com/twinj/uuid"
import "time"
//import "log"

const (
    CONN_HOST = "192.168.0.133"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func tcpListener(){
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
					fmt.Println("Error accepting: ", err.Error())
					os.Exit(1)
			}
			// Handle connections in a new goroutine.
			go handleRequest(conn)
	}
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


func main(){
	uuid.Init()
	fmt.Println("WebServer Test Connection")
	//for accepting TCP connections and executing a command
	//tcpListener();
	commandListener();
}


func commandExecuter(command string){
		result := processCommand(command)

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

		if result[0] == "SELL"{
			fmt.Println(len(result))
			sell(result[1],result[2],result[3])
		}

		if result[0] == "SELL"{
			fmt.Println(len(result))
			//sell(result[0],result[1],result[2],result[3])
		}
}

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
	}
}

func processCommand(text string) []string{
	result := strings.Split(text, ",")
	for i := range result {
		fmt.Println(result[i])
	}
	return result;
}

//depricated - no longer used
//manually select function to test
func selectCommand(text string){
	fmt.Print(text)
	if text == "quote\n"{
		fmt.Print("quoteTest")
		//quoteRequest()
	}
	if text == "adduser\n"{
		fmt.Print("addUser Test")
		//addUser()
	}

	if text == "buy\n"{
		fmt.Print("buy Test")
		//buy()
	}
	if text == "sell\n"{
		fmt.Println("sell Test")
		//sell()
	}
}

func quoteRequest(userId string, stockSymbol string) []string{
	// connect to this socket
	conn, _ := net.Dial("tcp", "192.168.3.104:3333")
		// read in input from stdin
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Text to send: ")
		//text, _ := reader.ReadString('\n')
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

		return messageArray
}

func createSession(){
	//create db session
	/*
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "userdb"
	session, _ := cluster.CreateSession()
	*/
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


func addUser(userId string, usableCashString string){

	usableCash := stringToCents(usableCashString)
	fmt.Println(usableCash)
	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	usableCashString = strconv.FormatInt(int64(usableCash), 10)

	//if err := session.Query("INSERT INTO users (userid, usableCash) VALUES ('Jones', 351) IF NOT EXISTS").Exec(); err != nil {
	//	panic(fmt.Sprintf("problem creating session", err))
	//}

	if err := session.Query("INSERT INTO users (userid, usableCash) VALUES ('" + userId + "', " + usableCashString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	defer session.Close()
}

//checks the state and runs only after a buy or sell to check if the UUID of a transaction is expired or NOT
//this is needed to return the allocated money in the case the transaction automatically expires
//waits for 62 seconds, checks the UUID parameter if it exists in the redis database, and if it doesnt it will revert the buy or sell command
func updateStateBuy(operation int, uuid string, userId string){

	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	timer1 := time.NewTimer(time.Second * 62)

	<-timer1.C
    fmt.Println("Timer1 has expired")
		fmt.Println("User Cash will be returned")

		if operation == 1 {
			fmt.Println("starting operation 1")
			var pendingCash int
			var usableCash int
			var totalCash int


			//obtain value remaining for expired transaction
			if err := session.Query("select pendingCash from buypendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&pendingCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			if pendingCash == 0{
				return
			}

			//delete pending transaction
			if err := session.Query("delete from buypendingtransactions where pid=" + uuid  + " and userid='" + userId + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}
			//obtain current users bank account value
			if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//add back accout value to pending cash
			totalCash = pendingCash + usableCash;
			totalCashString := strconv.FormatInt(int64(totalCash), 10)

			//return total money to user
			if err := session.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

		}
}

func updateStateSell(userId string, uuid string, usid string){
	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	print("In update sell")
	timer1 := time.NewTimer(time.Second * 62)

	<-timer1.C
    fmt.Println("Timer1 has expired")
	fmt.Println("User stocks will be returned")

	var pendingCash int
	var pendingStocks int
	var currentStocks int
	var totalStocks int
	fmt.Println(usid)
	fmt.Println(uuid)

	//obtain number of stocks for expired transaction
	if err := session.Query("select pendingcash, stockvalue from sellpendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&pendingCash, &pendingStocks); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := session.Query("delete from sellpendingtransactions where pid=" + uuid  + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	//get current users stock amount
	if err := session.Query("select stockamount from userstocks where usid="+usid).Scan(&currentStocks); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//add back stocks to stocks
	stocks := pendingCash/pendingStocks
	totalStocks = stocks + currentStocks
	totalStocksString := strconv.FormatInt(int64(totalStocks), 10)

	//return total stocks to the userstocks
	if err := session.Query("UPDATE userstocks SET stockamount =" + totalStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}

func cancelBuy(userId string){
	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	var pendingCash int
	var usableCash int
	var totalCash int
	var uuid string
	userId = strings.TrimSuffix(userId, "\n")

	//grab usableCash of the user
	if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//retrieve pending cash from most recent buy transaction
	if err := session.Query("select pid, pendingCash from buypendingtransactions where userId='" + userId + "'" + " LIMIT 1").Scan(&uuid, &pendingCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	totalCash = pendingCash + usableCash;
	totalCashString := strconv.FormatInt(int64(totalCash), 10)

	if err := session.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := session.Query("delete from buypendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}

func commitBuy(userId string){

	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	var pendingCash int
	var buyingstockName string
	var stockValue int
	var buyableStocks int
	var remainingCash int
	var usableCash int
	var uuid string
	userId = strings.TrimSuffix(userId, "\n")


	if err := session.Query("select pid, stock, stockValue, pendingCash from buypendingtransactions where userId='" + userId + "'").Scan(&uuid,&buyingstockName, &stockValue, &pendingCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	var usid string
	var ownedstockname string
	var stockamount int
	var hasStock bool

	//check if user currently owns any of this stock
	iter := session.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='"+ userId + "'").Iter()
	for iter.Scan(&usid, &ownedstockname, &stockamount) {
		if (ownedstockname == buyingstockName){
			hasStock = true
			break;
		}
		//fmt.Println("STOCKS: ", stockname, stockamount)
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if hasStock == true{

		//calculate amount of stocks can be bought
		buyableStocks = pendingCash / stockValue
		buyableStocks = buyableStocks + stockamount
		//remaining money
		remainingCash = pendingCash - (buyableStocks * stockValue)

		buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)

		//insert new stock record
		if err := session.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE usid=" + usid).Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//check users available cash
		if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//add available cash to leftover cash
		usableCash = usableCash + remainingCash
		usableCashString := strconv.FormatInt(int64(usableCash), 10)

		//re input the new cash value in to the user db
		if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	} else {
		//IF USE DOESNT OWN ANY OF THIS STOCK
		//calculate amount of stocks can be bought
		buyableStocks = pendingCash / stockValue
		//remaining money
		remainingCash = pendingCash - (buyableStocks * stockValue)

		buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)

		//insert new stock record
		if err := session.Query("INSERT INTO userstocks (usid, userid, stockamount, stockname) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + buyingstockName + "')").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//check users available cash
		if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//add available cash to leftover cash
		usableCash = usableCash + remainingCash
		usableCashString := strconv.FormatInt(int64(usableCash), 10)

		//re input the new cash value in to the user db
		if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}

	//delete the pending transcation
	if err := session.Query("delete from buypendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func buy(userId string, stock string, pendingCashString string){
	//userid,stocksymbol,amount

	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	pendingCash := stringToCents(pendingCashString)
	//var userId string = "Jones"
	//cash to spend in total for a stock
	//var pendingCash int = 200
	//var stock string = "abs"
	var stockValue int
	var usableCash int

	message := quoteRequest(userId, stock)
	fmt.Println(message[0])
	stockValueQuoteString := message[0]
	stockValue = stringToCents(stockValueQuoteString)


	//check if user has enough money for a BUY
	if err := session.Query("select userId, usableCash from users where userid='" + userId + "'").Scan(&userId, &usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("\n" + userId);
	fmt.Println("usableCash");
	fmt.Println(usableCash);
	fmt.Println("pendingCash");
	fmt.Println(pendingCash);
	//if not close the session
	if usableCash < pendingCash{
		session.Close()
		return
	}

	//if has enough cash subtract and set aside from db
	usableCash = usableCash - pendingCash;
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString = strconv.FormatInt(int64(pendingCash), 10)
	fmt.Println("Available Cash is greater than buy amount");
	if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("Cash allocated");



	//make a record of the new transaction

	u := uuid.NewV1()
	f := uuid.Formatter(u, uuid.FormatCanonical)
	fmt.Println(f)

	stockValueString := strconv.FormatInt(int64(stockValue), 10)
	if err := session.Query("INSERT INTO buypendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//---**************---------insert userid and pid into the redis database to start decrementing the transaction-------*********---------


	// NEED TO HAVE SMOETHING TO CHECK WHEN THE 60 seconds is up to return the money back and alert the user

	//run update function to check if the buy command has expired
	go updateStateBuy(1, f, userId);


	defer session.Close()
}

//sets aside the amount of money user wants to spend on a given stock
//executes prior to setTriggerValue
func setBuyAmount(userId string, stock string, pendingCashString string){

	//create session with cass database
	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//Verify that use funds is greater than amount attempting to spend


	//Usable Cash is stored as cents
	var usableCash int


	//convert pendingCash from string to int of cents
	pendingCash := stringToCents(pendingCashString)


	if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem getting usable cash form users", err))
	}

	//Verify the pending cash vs the usable cash
	fmt.Println("\n" + userId)
	fmt.Println("usableCash")
	fmt.Println(usableCash)
	fmt.Println("pendingCash")
	fmt.Println(pendingCash)

	//if the user doesnt have enough funds cancel the allocation
	if usableCash < pendingCash{
		fmt.Println("Not enough money for this transaction")
		session.Close()
		return
	}

	//allocate cash after being verified
	usableCash = usableCash - pendingCash;
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString = strconv.FormatInt(int64(pendingCash), 10)
	fmt.Println("Available Cash is greater than buy amount")
	if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem getting allocating user funds", err))
	}
	fmt.Println("Cash allocated")

	//Create an entry in the "Triggers" table to keep track of the initial buy amount setting

	//generate UUID to input as a unique value
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	//buy operation flag
	var operation string = "true"

	if err := session.Query("INSERT INTO triggers (tid, operation, pendingCash, stock, userid) VALUES (" + f + ", " + operation + ", " + pendingCashString + ", '" + stock + "', '" + userId + "')").Exec(); err != nil{
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}

	defer session.Close()

}

//Set maxmimum price of a stock before the stock gets auto bought
func setBuyTrigger(userId string, stock string, stockPriceTriggerString string){

	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//convert trigger price from string to int cents
	stockPriceTrigger := stringToCents(stockPriceTriggerString)
	fmt.Println(stockPriceTrigger);

	stockPriceTriggerString = strconv.FormatInt(int64(stockPriceTrigger), 10)

	//set the triggerValue and create thread to check the quote server

	if err := session.Query("UPDATE triggers SET triggerValue =" + stockPriceTriggerString + " WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem setting trigger", err))
	}

	go checkBuyTrigger(userId, stock, stockPriceTrigger)

}

func checkBuyTrigger(userId string, stock string, stockPriceTrigger int){

	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	for {
		//check the quote server every 5 seconds
		timer1 := time.NewTimer(time.Second * 5)
		<-timer1.C

		//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
		//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
		//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
		//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------

		message := quoteRequest(userId, stock)
		currentStockPrice := stringToCents(message[0])

		fmt.Println("Trigger value")
		fmt.Println(stockPriceTrigger)
		fmt.Println("quote price")
		fmt.Println(currentStockPrice)

		//execute the buy instantly if trigger condition is true
		if(currentStockPrice <= stockPriceTrigger){

			var usableCash int
			var pendingCash int
			stockValue := currentStockPrice
			var remainingCash int
			var usid string
			var ownedstockname string
			var stockamount int
			var hasStock bool

			//check if user currently owns any of this stock
			iter := session.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='"+ userId + "'").Iter()
			for iter.Scan(&usid, &ownedstockname, &stockamount) {
				if (ownedstockname == stock){
					hasStock = true
					break;
				}
			}
			if err := iter.Close(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}


			//If the user has some stock, add it to currently owned
			if hasStock == true{

				//grab pendingCash for the buy trigger
				if err := session.Query("SELECT pendingCash FROM triggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&pendingCash); err != nil {
					panic(fmt.Sprintf("problem getting usable cash form users", err))
				}

				//calculate amount of stocks can be bought
				buyableStocks := pendingCash / stockValue
				buyableStocks = buyableStocks + stockamount
				//remaining money
				remainingCash = pendingCash - (buyableStocks * stockValue)

				buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)


				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------





				//insert new stock record
				if err := session.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE usid=" + usid + "").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//check users available cash
				if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//add available cash to leftover cash
				usableCash = usableCash + remainingCash
				usableCashString := strconv.FormatInt(int64(usableCash), 10)

				//re input the new cash value in to the user db
				if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				if err := session.Query("DELETE FROM triggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				return

			} else {

				//get pending cash in the trigger
				if err := session.Query("SELECT pendingCash FROM triggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&pendingCash); err != nil {
					panic(fmt.Sprintf("problem getting usable cash form users", err))
				}

				//IF USE DOESNT OWN ANY OF THIS STOCK
				//calculate amount of stocks can be bought
				buyableStocks := pendingCash / stockValue
				fmt.Println("buyable stock amount")
				fmt.Println(buyableStocks)
				//remaining money
				remainingCash = pendingCash - (buyableStocks * stockValue)

				buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)

				fmt.Println("buyable stock string amount")
				fmt.Println(buyableStocksString)

				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------
				//CHECK IF TRIGGER IS STILL VALID AND HASNT BEEN DELETED WITHIN THE TIME OF EXECUTION----------------------------



				//insert new stock record
				if err := session.Query("INSERT INTO userstocks (usid, userid, stockamount, stockname) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + stock + "')").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//check users available cash
				if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//add available cash to leftover cash
				usableCash = usableCash + remainingCash
				usableCashString := strconv.FormatInt(int64(usableCash), 10)

				//re input the new cash value in to the user db
				if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				if err := session.Query("DELETE FROM triggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				return

			}
		}
	}
}

//cancel any buy triggers as well as buy_sell_amounts
func cancelBuyTrigger(userId string, stock string){

	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if err := session.Query("DELETE FROM triggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}

func sell(userId string, stock string, sellStockDollarsString string){
//userid,stocksymbol,amount
	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//var userId string = "Jones"
	sellStockDollars := stringToCents(sellStockDollarsString)
	//var stock string = "abc"
	var stockValue int
	var usableStocks int
	var stockname string
	var stockamount int
	var usid string
	var hasStock bool

	message := quoteRequest(userId, stock)
	fmt.Println(message[0])
	stockValue = stringToCents(message[0])

	//check if user has enough stocks for a SELL

	iter := session.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='Jones'").Iter()
	for iter.Scan(&usid, &stockname, &stockamount) {
		if (stockname == stock){
			hasStock = true
			break;

		}
		//fmt.Println("STOCKS: ", stockname, stockamount)
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	println(hasStock)
	if (!hasStock){
		session.Close()
		return
	}
	fmt.Println(stockname, stockamount)
	fmt.Println("\n" + userId);
	usableStocks = stockamount
	fmt.Println(usableStocks);

	//if not close the session
	if  (stockValue*usableStocks) < sellStockDollars{
		session.Close()
		return
	}
	var sellableStocks int = sellStockDollars/stockValue
	print("total sellable stocks ")
	print(sellableStocks)
	//if has enough stock for desired sell amount, set aside stocks from db
	usableStocks = usableStocks - sellableStocks;
	usableStocksString := strconv.FormatInt(int64(usableStocks),10)
	pendingCash := sellableStocks * stockValue;
	pendingCashString := strconv.FormatInt(int64(pendingCash), 10)
	stockValueString := strconv.FormatInt(int64(stockValue), 10)
	fmt.Println("Available Stocks is greater than sell amount");
	if err := session.Query("UPDATE userstocks SET stockamount =" + usableStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("Stocks allocated");

	u := uuid.NewV1()
	/*
	fmt.Println(id)
	u := uuid.NewV4()
	*/
	f := uuid.Formatter(u, uuid.FormatCanonical)
	fmt.Println(f)
	
	//tm := time.Now()

	//make a record of the new transaction

	if err := session.Query("INSERT INTO sellpendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	print("going to updateStateSell")
	go updateStateSell(userId, f, usid)

	defer session.Close()
}

func commitSell(userId string){
	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	var uuid string
	var pendingCash int
	var usableCash int
	userId = strings.TrimSuffix(userId, "\n")

	if err := session.Query("select pid from sellpendingtransactions where userid='" + userId + "'").Scan(&uuid); err != nil{
		panic(fmt.Sprintf("problem", err))
	}

	//get pending cash to be added to user account
	if err := session.Query("select pid, pendingcash from sellpendingtransactions where userid='" + userId + "'").Scan(&uuid, &pendingCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if uuid == "" {
		return
	}

	//get current users cash
	if err := session.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	//add available cash to leftover cash
	usableCash = usableCash + pendingCash
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	fmt.Println(usableCashString)

	//re input the new cash value in to the user db
	if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	//delete the pending transcation
	if err := session.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}

func cancelSell(userId string){
	cluster := gocql.NewCluster("192.168.3.103")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	var uuid string
	var pendingCash int
	var pendingStock int
	var stock string
	var usid string
	var stockname string
	var stockamount int
	var totalStocks int
	var stocks int

	userId = strings.TrimSuffix(userId, "\n")

	//obtain stocks and number of stocks for cancelled transaction
	if err := session.Query("select pid, userId, pendingcash, stock, stockvalue from sellpendingtransactions where userId='" + userId + "'" + " LIMIT 1").Scan(&uuid, &userId, &pendingCash, &stock, &pendingStock, ); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//get current user stocks
	iter := session.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='Jones'").Iter()
	for iter.Scan(&usid, &stockname, &stockamount) {
		if (stockname == stock){
			break;

		}
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//convert stock value to stock amount
	stocks = pendingCash/pendingStock
	totalStocks =  stocks + stockamount
	totalStocksString := strconv.FormatInt(int64(totalStocks),10)

	//return user stocks
	if err := session.Query("UPDATE userstocks SET stockamount =" + totalStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := session.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}

func deleteSession(){

}

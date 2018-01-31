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
    CONN_HOST = "localhost"
    CONN_PORT = "3334"
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

var sessionGlobal *gocql.Session
var transactionNumGlobal int

func main(){

	cluster := gocql.NewCluster("172.17.0.3")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	sessionGlobal = session
	
	transactionNumGlobal = 0;
	

	uuid.Init()
	fmt.Println("WebServer Test Connection")
	//for accepting TCP connections and executing a command
	tcpListener();
	//commandListener();
}


func commandExecuter(command string){
		result := processCommand(command)

		if result[0] == "ADD"{
			transactionNumGlobal += 1
			go addUser(result[1],result[2],transactionNumGlobal)
		}

		if result[0] == "QUOTE"{
			transactionNumGlobal += 1
			go quoteRequest(result[1],result[2],transactionNumGlobal)
		}

		if result[0] == "BUY"{
			transactionNumGlobal += 1
			fmt.Println(result[1])
			fmt.Println(result[2])
			fmt.Println(result[3])
			go buy(result[1],result[2],result[3],transactionNumGlobal)
		}

		if result[0] == "COMMIT_BUY"{
			transactionNumGlobal += 1
			go commitBuy(result[1],transactionNumGlobal)
		}

		if result[0] == "CANCEL_BUY"{
			transactionNumGlobal += 1
			go cancelBuy(result[1],transactionNumGlobal)
		}

		if result[0] == "SELL"{
			transactionNumGlobal += 1
			go sell(result[1],result[2],result[3],transactionNumGlobal)
		}

		if result[0] == "COMMIT_SELL"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go commitSell(result[1],transactionNumGlobal)
		}

		if result[0] == "CANCEL_SELL"{
			transactionNumGlobal += 1
			go cancelSell(result[1],transactionNumGlobal)
		}

		if result[0] == "SET_BUY_AMOUNT"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go setBuyAmount(result[1],result[2],result[3],transactionNumGlobal)
		}

		if result[0] == "SET_BUY_TRIGGER"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go setBuyTrigger(result[1],result[2],result[3],transactionNumGlobal)
		}

		if result[0] == "SET_SELL_AMOUNT"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go setSellAmount(result[1],result[2],result[3],transactionNumGlobal)
		}

		if result[0] == "SET_SELL_TRIGGER"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go setSellTrigger(result[1],result[2],result[3],transactionNumGlobal)
		}
		if result[0] == "CANCEL_SET_SELL"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go cancelSellTrigger(result[1],result[2],transactionNumGlobal)
		}
		if result[0] == "CANCEL_SET_BUY"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go cancelBuyTrigger(result[1],result[2],transactionNumGlobal)
		}

		if result[0] == "DISPLAY_SUMMARY"{
			transactionNumGlobal += 1
			fmt.Println(len(result))
			go displaySummary(result[1])
		}

		if result[0] == "DUMPLOG"{
			transactionNumGlobal += 1
			if len(result) == 3{
				go dumpUser(result[1],result[2])
			} else if len(result) == 2{
				go dump(result[1])
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

func displaySummary(userId string){
	//return user summary of their stocks, cash, triggers, etc
}

func quoteRequest(userId string, stockSymbol string ,transactionNum int) []string{
	// connect to this socket
		//conn, _ := net.Dial("tcp", "quoteserve.seng:3333")
		conn, _ := net.Dial("tcp", "quoteserve.seng:4446")
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

		timestamp_q := (time.Now().UTC().UnixNano())/ 1000000
		timestamp_quote := strconv.FormatInt(timestamp_q,10)
		transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
		go logQuoteEvent(timestamp_quote,"TS1",transactionNum_string,messageArray[0],messageArray[1],userId,messageArray[3],messageArray[4])

		return messageArray
}

func logUserEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string){
	//connect to audit server
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds
	fmt.Fprintf(conn,text + "\n")
	//go logSystemEvent(time, "AU1", "1",command,userid,"","")
}

func logQuoteEvent(time string, server string, transactionNum string, price string, stockSymbol string, userid string, quoteservertime string, cryptokey string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "Quote" + "," + time + "," + server + "," + transactionNum + "," + price + "," + stockSymbol + "," + userid + "," + quoteservertime + "," + cryptokey
	fmt.Fprintf(conn,text + "\n") 
}

func logSystemEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "System" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds
	fmt.Fprintf(conn,text + "\n") 
}

func logAccountTransactionEvent(time string, server string, transactionNum string, action string, userid string, funds string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "Account" + "," + time + "," + server + "," + transactionNum + "," + action + "," + userid + "," + funds
	fmt.Fprintf(conn,text + "\n") 	
}

func logErrorEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, error string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds + "," + error
	fmt.Fprintf(conn,text + "\n") 
}

func logDebugEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, debug string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "User" + "," + time + "," + server + "," + transactionNum + "," + command + "," + userid + "," + stockSymbol + "," + funds + "," + debug
	fmt.Fprintf(conn,text + "\n") 
}

func dumpUser(userId string, filename string){
	fmt.Println("In Dump user")
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "DUMPLOG" + "," + userId + "," + filename
	fmt.Fprintf(conn,text + "\n") 
}

func dump(filename string){
	conn, _ := net.Dial("tcp", "localhost:5555")
	text := "DUMPLOG" + "," + filename
	fmt.Fprintf(conn,text + "\n") 
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


func addUser(userId string, usableCashString string,transactionNum int){


	usableCash := stringToCents(usableCashString)

	fmt.Println(usableCash)

	var count int

	if err := sessionGlobal.Query("SELECT count(*) FROM users WHERE userid='" + userId + "'").Scan(&count); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//if the user already exists add money to the account
	if count != 0{
		fmt.Println("adding funds to user")
		addFunds(userId, usableCash)


	//if the user doesnt exist create a new user
	}else{
		fmt.Println("creating new user")
		usableCashString = strconv.FormatInt(int64(usableCash), 10)
		if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('" + userId + "', " + usableCashString + ")").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}	
	}

	//usableCashString = strconv.FormatInt(int64(usableCash), 10)

	//if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('Jones', 351) IF NOT EXISTS").Exec(); err != nil {
	//	panic(fmt.Sprintf("problem creating session", err))
	//}

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logAccountTransactionEvent(timestamp_command, "TS1", "1", "ADD", userId, usableCashString)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "ADD", userId, "", usableCashString)

	//if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('" + userId + "', " + usableCashString + ")").Exec(); err != nil {
	//	panic(fmt.Sprintf("problem creating session", err))
	//}

	
}

//checks the state and runs only after a buy or sell to check if the UUID of a transaction is expired or NOT
//this is needed to return the allocated money in the case the transaction automatically expires
//waits for 62 seconds, checks the UUID parameter if it exists in the redis database, and if it doesnt it will revert the buy or sell command
func updateStateBuy(operation int, uuid string, userId string){


	timer1 := time.NewTimer(time.Second * 62)

	<-timer1.C
    fmt.Println("Timer1 has expired")
	fmt.Println("User Cash will be returned")


		if operation == 1 {
			fmt.Println("starting operation 1")
			var pendingCash int
			var usableCash int
			var totalCash int
			var count int

			//check if remaining transaction still exists
			fmt.Println("Checking if the the buy transaction still exists")
			if err := sessionGlobal.Query("select count(*) from buypendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&count); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			fmt.Println("pending buy transactions:")
			fmt.Println(count)
			if count == 0 {
				fmt.Println("buy transaction doesnt exist")
				return;
			}

			//obtain value remaining for expired transaction
			if err := sessionGlobal.Query("select pendingCash from buypendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&pendingCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			if pendingCash == 0{
				return
			}

			//delete pending transaction
			if err := sessionGlobal.Query("delete from buypendingtransactions where pid=" + uuid  + " and userid='" + userId + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}
			//obtain current users bank account value
			if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//add back accout value to pending cash
			totalCash = pendingCash + usableCash;
			totalCashString := strconv.FormatInt(int64(totalCash), 10)

			//return total money to user
			if err := sessionGlobal.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

		}
}

func updateStateSell(userId string, uuid string, usid string){
	print("In update sell")
	timer1 := time.NewTimer(time.Second * 62)

	<-timer1.C
    fmt.Println("Timer1 has expired for SELL command")
	fmt.Println("User stocks will be returned")

	var pendingCash int
	var pendingStocks int
	var currentStocks int
	var totalStocks int
	//fmt.Println(usid)
	//fmt.Println(uuid)
	var count int

	//check if remaining transaction still exists
	if err := sessionGlobal.Query("select count(*) from sellpendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&count); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if count == 0 {
		return;
	}

	//obtain number of stocks for expired transaction
	if err := sessionGlobal.Query("select pendingcash, stockvalue from sellpendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&pendingCash, &pendingStocks); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := sessionGlobal.Query("delete from sellpendingtransactions where pid=" + uuid  + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	//get current users stock amount
	if err := sessionGlobal.Query("select stockamount from userstocks where usid="+usid).Scan(&currentStocks); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//add back stocks to stocks
	stocks := pendingCash/pendingStocks
	totalStocks = stocks + currentStocks
	totalStocksString := strconv.FormatInt(int64(totalStocks), 10)

	//return total stocks to the userstocks
	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + totalStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}

func cancelBuy(userId string,transactionNum int){
	var pendingCash int
	var usableCash int
	var totalCash int
	var uuid string
	var stock string
	userId = strings.TrimSuffix(userId, "\n")

	sellExists := checkDependency("CANCEL_BUY",userId,"none")
	if(sellExists == false){
		fmt.Println("cannot CANCEL BUY, no buys pending")
		return
	}

	//grab usableCash of the user
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//retrieve pending cash from most recent buy transaction
	if err := sessionGlobal.Query("select pid, pendingCash, stock from buypendingtransactions where userId='" + userId + "'" + " LIMIT 1").Scan(&uuid, &pendingCash,&stock); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	totalCash = pendingCash + usableCash;
	totalCashString := strconv.FormatInt(int64(totalCash), 10)

	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := sessionGlobal.Query("delete from buypendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "CANCEL_BUY", userId, stock, "")

}

func commitBuy(userId string,transactionNum int){


	buyExists := checkDependency("COMMIT_BUY",userId,"none")
	if(buyExists == false){
		fmt.Println("cannot commit, no buys pending")
		return
	}

	var pendingCash int
	var buyingstockName string
	var stockValue int
	var buyableStocks int
	var remainingCash int
	var usableCash int
	var uuid string
	userId = strings.TrimSuffix(userId, "\n")


	if err := sessionGlobal.Query("select pid, stock, stockValue, pendingCash from buypendingtransactions where userId='" + userId + "'").Scan(&uuid,&buyingstockName, &stockValue, &pendingCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	var usid string
	var ownedstockname string
	var stockamount int
	var hasStock bool

	//check if user currently owns any of this stock
	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='"+ userId + "'").Iter()
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
		if err := sessionGlobal.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE usid=" + usid).Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//check users available cash
		if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//add available cash to leftover cash
		usableCash = usableCash + remainingCash
		usableCashString := strconv.FormatInt(int64(usableCash), 10)

		//re input the new cash value in to the user db
		if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
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
		if err := sessionGlobal.Query("INSERT INTO userstocks (usid, userid, stockamount, stockname) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + buyingstockName + "')").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//check users available cash
		if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//add available cash to leftover cash
		usableCash = usableCash + remainingCash
		usableCashString := strconv.FormatInt(int64(usableCash), 10)

		//re input the new cash value in to the user db
		if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "COMMIT_BUY", userId, buyingstockName, "")

	//delete the pending transcation
	if err := sessionGlobal.Query("delete from buypendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


}

func buy(userId string, stock string, pendingCashString string,transactionNum int){
	//userid,stocksymbol,amount


	pendingCash := stringToCents(pendingCashString)
	//var userId string = "Jones"
	//cash to spend in total for a stock
	//var pendingCash int = 200
	//var stock string = "abs"
	var stockValue int
	var usableCash int

	message := quoteRequest(userId, stock,transactionNum)

	timestamp_q := (time.Now().UTC().UnixNano())/ 1000000
	timestamp_quote := strconv.FormatInt(timestamp_q,10)
	//transactionNum_quote += 1
	//transactionNum_quote_string := strconv.FormatInt(int64(transactionNum_quote), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logQuoteEvent(timestamp_quote,"TS1",transactionNum_string,message[0],message[1],userId,message[3],message[4])
	
	fmt.Println(message[0])
	stockValueQuoteString := message[0]
	stockValue = stringToCents(stockValueQuoteString)


	//check if user has enough money for a BUY
	if err := sessionGlobal.Query("select userId, usableCash from users where userid='" + userId + "'").Scan(&userId, &usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("\n" + userId);
	fmt.Println("usableCash");
	fmt.Println(usableCash);
	fmt.Println("pendingCash");
	fmt.Println(pendingCash);
	//if not close the session
	if usableCash < pendingCash{
		
		return
	}

	//if has enough cash subtract and set aside from db
	usableCash = usableCash - pendingCash;
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString = strconv.FormatInt(int64(pendingCash), 10)
	fmt.Println("Available Cash is greater than buy amount");
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("Cash allocated");



	//make a record of the new transaction

	u := uuid.NewV1()
	f := uuid.Formatter(u, uuid.FormatCanonical)
	fmt.Println(f)


	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "BUY", userId, stock, pendingCashString)

	
	stockValueString := strconv.FormatInt(int64(stockValue), 10)
	if err := sessionGlobal.Query("INSERT INTO buypendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//---**************---------insert userid and pid into the redis database to start decrementing the transaction-------*********---------


	// NEED TO HAVE SMOETHING TO CHECK WHEN THE 60 seconds is up to return the money back and alert the user

	//run update function to check if the buy command has expired
	go updateStateBuy(1, f, userId);


	
}

//sets aside the amount of money user wants to spend on a given stock
//executes prior to setTriggerValue
func setBuyAmount(userId string, stock string, pendingCashString string,transactionNum int){

	//create session with cass database

	//Verify that use funds is greater than amount attempting to spend


	//Usable Cash is stored as cents
	var usableCash int


	//convert pendingCash from string to int of cents
	pendingCash := stringToCents(pendingCashString)


	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
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
		
		return
	}

	//allocate cash after being verified
	usableCash = usableCash - pendingCash;
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString = strconv.FormatInt(int64(pendingCash), 10)
	fmt.Println("Available Cash is greater than buy amount")
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem getting allocating user funds", err))
	}
	fmt.Println("Cash allocated")

	//Create an entry in the "Triggers" table to keep track of the initial buy amount setting

	//generate UUID to input as a unique value
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	//buy operation flag
	//var operation string = "true"

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "SET_BUY_AMOUNT", userId, stock, pendingCashString)

	if err := sessionGlobal.Query("INSERT INTO buyTriggers (tid, pendingCash, stock, userid) VALUES (" + f + ", " + pendingCashString + ", '" + stock + "', '" + userId + "')").Exec(); err != nil{
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	

}

//Set maxmimum price of a stock before the stock gets auto bought
func setBuyTrigger(userId string, stock string, stockPriceTriggerString string,transactionNum int){


	//convert trigger price from string to int cents
	stockPriceTrigger := stringToCents(stockPriceTriggerString)
	fmt.Println(stockPriceTrigger);

	stockPriceTriggerString = strconv.FormatInt(int64(stockPriceTrigger), 10)

	//set the triggerValue and create thread to check the quote server

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "SET_BUY_TRIGGER", userId, stock, stockPriceTriggerString)

	if err := sessionGlobal.Query("UPDATE buyTriggers SET triggerValue =" + stockPriceTriggerString + " WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem setting trigger", err))
	}


	go checkBuyTrigger(userId, stock, stockPriceTrigger,transactionNum)

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

func checkBuyTrigger(userId string, stock string, stockPriceTrigger int,transactionNum int){



	operation := true

	for {
		//check the quote server every 5 seconds
		timer1 := time.NewTimer(time.Second * 5)
		<-timer1.C

		//if the trigger doesnt exist exit

		exists := checkTriggerExists(userId, stock, operation)
		if exists == false{
			return
		}



		message := quoteRequest(userId, stock,transactionNum)
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
			iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='"+ userId + "'").Iter()
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
				if err := sessionGlobal.Query("SELECT pendingCash FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&pendingCash); err != nil {
					panic(fmt.Sprintf("problem getting usable cash form users", err))
				}

				//calculate amount of stocks can be bought
				buyableStocks := pendingCash / stockValue
				buyableStocks = buyableStocks + stockamount
				//remaining money
				remainingCash = pendingCash - (buyableStocks * stockValue)

				buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)


				//if the trigger doesnt exist exit
				exists := checkTriggerExists(userId, stock, operation)
				if exists == false{
					return
				}


				//insert new stock record
				if err := sessionGlobal.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE usid=" + usid + "").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//check users available cash
				if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//add available cash to leftover cash
				usableCash = usableCash + remainingCash
				usableCashString := strconv.FormatInt(int64(usableCash), 10)

				//re input the new cash value in to the user db
				if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				if err := sessionGlobal.Query("DELETE FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				return

			} else {

				//get pending cash in the trigger
				if err := sessionGlobal.Query("SELECT pendingCash FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&pendingCash); err != nil {
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

				//if the trigger doesnt exist exit
				exists := checkTriggerExists(userId, stock, operation)
				if exists == false{
					return
				}


				//insert new stock record
				if err := sessionGlobal.Query("INSERT INTO userstocks (usid, userid, stockamount, stockname) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + stock + "')").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//check users available cash
				if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//add available cash to leftover cash
				usableCash = usableCash + remainingCash
				usableCashString := strconv.FormatInt(int64(usableCash), 10)

				//re input the new cash value in to the user db
				if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				if err := sessionGlobal.Query("DELETE FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				return

			}
		}
	}
}

//cancel any buy triggers as well as buy_sell_amounts
func cancelBuyTrigger(userId string, stock string,transactionNum int){


	buyExists := checkDependency("CANCEL_SET_BUY",userId,stock)
	if(buyExists == false){
		fmt.Println("cannot CANCEL, no buys pending")
		return
	}

	fmt.Println("cancelling buy trigger")

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "CANCEL_SET_BUY", userId, stock, "")

	if err := sessionGlobal.Query("DELETE FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}

//cancels any sell triggers or sell amounts
func cancelSellTrigger(userId string, stock string,transactionNum int){


	sellExists := checkDependency("CANCEL_SET_SELL",userId,stock)
	if(sellExists == false){
		fmt.Println("cannot CANCEL, no sells pending")
		return
	}

	fmt.Println("cancelling sell trigger")

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "CANCEL_SET_SELL", userId, stock, "")

	if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}

//sets the total cash to gain from selling a stock
func setSellAmount(userId string, stock string, pendingCashString string,transactionNum int){


	pendingCashCents := stringToCents(pendingCashString)
	//check if user owns stock
	ownedStockAmount, usid := checkStockOwnership(userId, stock)
	fmt.Println(usid)

	if(ownedStockAmount == 0){
		fmt.Println("Cannot Sell a stock you don't own")
		return
	}

	pendingCashString = strconv.FormatInt(int64(pendingCashCents), 10)

	//create trigger to sell a certain amount of the stock
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "SET_SELL_AMOUNT", userId, stock, pendingCashString)

	//Create new entry for the sell trigger with the sell amount
	if err := sessionGlobal.Query("INSERT INTO sellTriggers (tid, pendingCash, stock, userid) VALUES (" + f + ", " + pendingCashString + ", '" + stock + "', '" + userId + "')").Exec(); err != nil{
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}
}

func setSellTrigger(userId string, stock string, stockSellPrice string,transactionNum int){



	stockSellPriceCents := stringToCents(stockSellPrice)
	stockSellPriceCentsString := strconv.FormatInt(int64(stockSellPriceCents), 10)

	//check if set sell amount is set for this particular stock
	var count int
	if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId +"' AND stock='" + stock + "' ").Scan(&count); err != nil{
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}

	//if set sell amount isnt set return
	if count == 0 {
		fmt.Println("No set Sell amount placed")
		return
	}

	//update database entry with trigger value
	if err := sessionGlobal.Query("UPDATE sellTriggers SET triggerValue=" + stockSellPriceCentsString + " WHERE userid='" + userId +"' AND stock='" + stock + "' ").Exec(); err != nil{
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "SET_SELL_TRIGGER", userId, stock, stockSellPrice)
	go checkSellTrigger(userId, stock, stockSellPriceCents,transactionNum)
}


func checkSellTrigger(userId string, stock string, stockSellPriceCents int,transactionNum int){

	operation := false


	for {
		//check the quote server every 5 seconds
		timer1 := time.NewTimer(time.Second * 5)
		<-timer1.C

		//if the trigger doesnt exist exit

		exists := checkTriggerExists(userId, stock, operation)
		if exists == false{
			return
		}

		//retrieve current stock price
		message := quoteRequest(userId, stock,transactionNum)
		currentStockPrice := stringToCents(message[0])

		if currentStockPrice > stockSellPriceCents{

			//Check how many stocks the user can sell

			var pendingCash int

			if err := sessionGlobal.Query("SELECT pendingCash FROM sellTriggers WHERE userid='" + userId +"' AND stock='" + stock + "' ").Scan(&pendingCash); err != nil{
				panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
			}

			//calculate amount of stocks can be sold
			sellAbleStocksMax := pendingCash / currentStockPrice

			//check how many stocks the user owns
			ownedStocks, usid := checkStockOwnership(userId, stock)

			//check how many stocks the user can sell
			var sellAbleStocks int
			var remainingStocks int

			//check if user has more owned stocks than able to sell
			if sellAbleStocksMax < ownedStocks{
				sellAbleStocks = sellAbleStocksMax
				remainingStocks = ownedStocks - sellAbleStocks

				//calculate money gained from stocks selling
				sellAbleStockPrice := pendingCash - (sellAbleStocksMax * currentStockPrice)

				remainingStocksString := strconv.FormatInt(int64(remainingStocks),10)

				//update userStock database with new about of stock
				if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + remainingStocksString + " WHERE usid=" + usid ).Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//increase money in userAccount
				addFunds(userId, sellAbleStockPrice)

				//delete trigger

				 

				if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
			}else{
				//Case where user does not own enough stocks to sell the maximum amount
				//user must sell the most it can
				sellAbleStocks = ownedStocks

				//ownedStocksString, err := strconv.Atoi("ownedStocks")

				sellAbleStockPrice := pendingCash - (sellAbleStocks * currentStockPrice)
				remainingCash := pendingCash - sellAbleStockPrice
				sellAbleStockPrice = sellAbleStockPrice + remainingCash

				addFunds(userId, sellAbleStockPrice)

				if err := sessionGlobal.Query("DELETE FROM userstocks WHERE usid=" + usid ).Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//delete trigger
				if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId=" + userId + " AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
			}
		}

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


func sell(userId string, stock string, sellStockDollarsString string,transactionNum int){
//userid,stocksymbol,amount

	//var userId string = "Jones"
	sellStockDollars := stringToCents(sellStockDollarsString)
	//var stock string = "abc"
	var stockValue int
	var usableStocks int
	var stockname string
	var stockamount int
	var usid string
	var hasStock bool

	message := quoteRequest(userId, stock,transactionNum)


	timestamp_q := (time.Now().UTC().UnixNano())/ 1000000
	timestamp_quote := strconv.FormatInt(timestamp_q,10) 
	//transactionNum_quote += 1
	//transactionNum_quote_string := strconv.FormatInt(int64(transactionNum_quote), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logQuoteEvent(timestamp_quote,"TS1",transactionNum_string,message[0],message[1],userId,message[3],message[4])
	fmt.Println(message[0])
	stockValue = stringToCents(message[0])

	//check if user has enough stocks for a SELL

	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
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
		
		return
	}
	fmt.Println(stockname, stockamount)
	fmt.Println("\n" + userId);
	usableStocks = stockamount
	fmt.Println(usableStocks);

	//if not close the session
	if  (stockValue*usableStocks) < sellStockDollars{
		
		return
	}
	sellableStocks := sellStockDollars/stockValue
	print("total sellable stocks ")
	fmt.Println(sellableStocks)
	//if has enough stock for desired sell amount, set aside stocks from db
	usableStocks = usableStocks - sellableStocks;
	usableStocksString := strconv.FormatInt(int64(usableStocks),10)
	fmt.Println("set stocks to " + usableStocksString)
	pendingCash := sellableStocks * stockValue;
	pendingCashString := strconv.FormatInt(int64(pendingCash), 10)
	stockValueString := strconv.FormatInt(int64(stockValue), 10)
	fmt.Println("Available Stocks is greater than sell amount");
	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + usableStocksString + " WHERE usid=" + usid).Exec(); err != nil {
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
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	//transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "SELL", userId, stock, sellStockDollarsString)

	if err := sessionGlobal.Query("INSERT INTO sellpendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	
	go updateStateSell(userId, f, usid)

	
}

//-----------------------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------------------
//----------------------commit sell function needs to subtract the sold stocks from the userStocks row-------------------------
//-----------------------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------------------
func commitSell(userId string,transactionNum int){

	var uuid string
	var pendingCash int
	var usableCash int
	var stock string
	userId = strings.TrimSuffix(userId, "\n")

	sellExists := checkDependency("COMMIT_SELL",userId,"none")
	if(sellExists == false){
		fmt.Println("cannot commit, no sells pending")
		return
	}


	if err := sessionGlobal.Query("select pid,stock from sellpendingtransactions where userid='" + userId + "'").Scan(&uuid,&stock); err != nil{
		panic(fmt.Sprintf("problem", err))
	}

	//get pending cash to be added to user account
	if err := sessionGlobal.Query("select pid, pendingcash from sellpendingtransactions where userid='" + userId + "'").Scan(&uuid, &pendingCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if uuid == "" {
		return
	}

	//get current users cash
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	//add available cash to leftover cash
	usableCash = usableCash + pendingCash
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	fmt.Println(usableCashString)

	//re input the new cash value in to the user db
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	//subtract sold stocks from users owned stocks



	//delete the pending transcation
	if err := sessionGlobal.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "COMMIT_SELL", userId, stock, "")

}

func cancelSell(userId string,transactionNum int){

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

	sellExists := checkDependency("CANCEL_SELL",userId,"none")
	if(sellExists == false){
		fmt.Println("cannot CANCEL SELL, no sell pending")
		return
	}

	//obtain stocks and number of stocks for cancelled transaction
	if err := sessionGlobal.Query("select pid, userId, pendingcash, stock, stockvalue from sellpendingtransactions where userId='" + userId + "'" + " LIMIT 1").Scan(&uuid, &userId, &pendingCash, &stock, &pendingStock, ); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//get current user stocks
	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
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
	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + totalStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := sessionGlobal.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	go logUserEvent(timestamp_command, "TS1", transactionNum_string, "CANCEL_SELL", userId, stock, "")
}

func deleteSession(){

}

package main

import "net"
import "fmt"
import "bufio"
import "os"
import "github.com/gocql/gocql"
import "strconv"
import "github.com/go-redis/redis"
import "github.com/twinj/uuid"
//import "log"



func main(){
	fmt.Println("WebServer Test Connection")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	selectCommand(text);
}

func selectCommand(text string){
	fmt.Print(text)
	if text == "quote\n"{
		fmt.Print("quoteTest")
		quoteRequest()
	}
	if text == "adduser\n"{
		fmt.Print("addUser Test")
		addUser()
	}

	if text == "buy\n"{
		fmt.Print("buy Test")
		buy()
	}
	if text == "sell\n"{
		fmt.Println("sell Test")
		sell()
	}
}

func quoteRequest(){
	// connect to this socket
	conn, _ := net.Dial("tcp", "192.168.0.133:3333")
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text + "\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
}

func createSession(){
	//create db session
	/*
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "userdb"
	session, _ := cluster.CreateSession()
	*/
}


func addUser(){
	cluster := gocql.NewCluster("192.168.1.131")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


	//if err := session.Query("INSERT INTO users (userid, usableCash) VALUES ('Jones', 351) IF NOT EXISTS").Exec(); err != nil {
	//	panic(fmt.Sprintf("problem creating session", err))
	//}
	if err := session.Query("INSERT INTO users (userid, usableCash) VALUES ('Jones', 9000)").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	defer session.Close()
}

func buy(){
	//userid,stocksymbol,amount
	cluster := gocql.NewCluster("192.168.1.131")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	var userId string = "Jones"
	//cash to spend in total for a stock
	var pendingCash int = 200
	var stock string = "abs"
	var operation string = "true"
	var committed string = "false"
	var stockValue int = 20
	var usableCash int

	//check if user has enough money for a BUY
	if err := session.Query("select userId, usableCash from users where userid='" + userId + "'").Scan(&userId, &usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("\n" + userId);
	fmt.Println(usableCash);
	//if not close the session
	if usableCash < pendingCash{
		session.Close()
		return
	}

	//if has enough cash subtract and set aside from db
	usableCash = usableCash - pendingCash;
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString := strconv.FormatInt(int64(pendingCash), 10)
	fmt.Println("Available Cash is greater than buy amount");
	if err := session.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("Cash allocated");



	//make a record of the new transaction

	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	fmt.Println(f)

	stockValueString := strconv.FormatInt(int64(stockValue), 10)
	if err := session.Query("INSERT INTO pendingtransactions (pid, committed, operation, userid, pendingCash, stock, stockValue) VALUES (" + f + ", " + committed + ", " + operation + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//---**************---------insert userid and pid into the redis database to start decrementing the transaction-------*********---------

	client := redis.NewClient(&redis.Options{
    Network:  "tcp",
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
	})

	//set a 60 second TTL in redis for this key value
	err = client.Set(f, "1", 60000000000).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(f).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(f, val)


	defer session.Close()
}

func sell(){
//userid,stocksymbol,amount
	cluster := gocql.NewCluster("192.168.0.111")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	var userId string = "Jones"
	var sellStockDollars int = 200
	var stock string = "abc"
	var operation string = "true"
	var committed string = "false"
	var stockValue int = 20
	var usableStocks int
	var stockname string
	var stockamount int
	var usid string
	var hasStock bool

	//check if user has enough stocks for a SELL
	
	if err := session.Query("select userId from users where userid='" + userId + "'").Scan(&userId); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	/*
	if err := session.Query("INSERT INTO userstocks (usid, userid, stockamount, stockname) VALUES (uuid(), 'Jones', 20, 'abc')").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	*/
	
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

	

	//make a record of the new transaction

	if err := session.Query("INSERT INTO pendingtransactions (pid, committed, operation, userid, pendingCash, stock, stockValue) VALUES (uuid(), " + committed + ", " + operation + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


	defer session.Close()
}

func deleteSession(){

}

package main

import "net"
import "fmt"
import "bufio"
import "os"
import "github.com/gocql/gocql"
import "strconv"
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
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "userdb"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


	//if err := session.Query("INSERT INTO users (userid, usableCash) VALUES ('Jones', 351) IF NOT EXISTS").Exec(); err != nil {
	//	panic(fmt.Sprintf("problem creating session", err))
	//}
	if err := session.Query("INSERT INTO users (userid, usableCash) VALUES ('Jones', 351)").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	defer session.Close()
}

func buy(){
	//userid,stocksymbol,amount
	cluster := gocql.NewCluster("127.0.0.1")
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

	if err := session.Query("INSERT INTO pendingtransactions (uid, committed, operation, userid, pendingCash, stock, stockValue) VALUES (uuid(), " + committed + ", " + operation + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValue + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}


	defer session.Close()
}

func deleteSession(){

}

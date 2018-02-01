package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	//"strings"
	"time"

	"github.com/twinj/uuid"
	//"github.com/go-redis/redis"
	//"log"
)

func buy(userId string, stock string, pendingCashString string, transactionNum int) {
	//userid,stocksymbol,amount

	pendingCash := stringToCents(pendingCashString)
	//var userId string = "Jones"
	//cash to spend in total for a stock
	//var pendingCash int = 200
	//var stock string = "abs"
	var stockValue int
	var usableCash int

	message := quoteRequest(userId, stock, transactionNum)

	timestamp_q := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_quote := strconv.FormatInt(timestamp_q, 10)
	//transactionNum_quote += 1
	//transactionNum_quote_string := strconv.FormatInt(int64(transactionNum_quote), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum), 10)
	logQuoteEvent(timestamp_quote, "TS1", transactionNum_string, message[0], message[1], userId, message[3], message[4])

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "BUY", userId, stock, pendingCashString)

	fmt.Println(message[0])
	stockValueQuoteString := message[0]
	stockValue = stringToCents(stockValueQuoteString)

	//check if user has enough money for a BUY
	if err := sessionGlobal.Query("select userId, usableCash from users where userid='"+userId+"'").Scan(&userId, &usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("\n" + userId)
	fmt.Println("usableCash")
	fmt.Println(usableCash)
	fmt.Println("pendingCash")
	fmt.Println(pendingCash)
	//if not close the session
	if usableCash < pendingCash {

		return
	}

	//if has enough cash subtract and set aside from db
	usableCash = usableCash - pendingCash
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString = strconv.FormatInt(int64(pendingCash), 10)
	fmt.Println("Available Cash is greater than buy amount")
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	fmt.Println("Cash allocated")

	u := uuid.NewV1()
	f := uuid.Formatter(u, uuid.FormatCanonical)
	fmt.Println(f)

	stockValueString := strconv.FormatInt(int64(stockValue), 10)
	if err := sessionGlobal.Query("INSERT INTO buypendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//---**************---------insert userid and pid into the redis database to start decrementing the transaction-------*********---------

	// NEED TO HAVE SMOETHING TO CHECK WHEN THE 60 seconds is up to return the money back and alert the user

	//run update function to check if the buy command has expired
	updateStateBuy(1, f, userId)
}

//checks the state and runs only after a buy or sell to check if the UUID of a transaction is expired or NOT
//this is needed to return the allocated money in the case the transaction automatically expires
//waits for 62 seconds, checks the UUID parameter if it exists in the redis database, and if it doesnt it will revert the buy or sell command
func updateStateBuy(operation int, uuid string, userId string) {

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
			return
		}

		//obtain value remaining for expired transaction
		if err := sessionGlobal.Query("select pendingCash from buypendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&pendingCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		if pendingCash == 0 {
			return
		}

		//delete pending transaction
		if err := sessionGlobal.Query("delete from buypendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
		//obtain current users bank account value
		if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		//add back accout value to pending cash
		totalCash = pendingCash + usableCash
		totalCashString := strconv.FormatInt(int64(totalCash), 10)

		//return total money to user
		if err := sessionGlobal.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}
}

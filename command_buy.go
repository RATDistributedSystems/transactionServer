package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/twinj/uuid"
)

func buy(userId string, stock string, pendingCashString string, transactionNum int) {

	logUserEvent("TS1", transactionNum, "BUY", userId, stock, pendingCashString)
	pendingCash := stringToCents(pendingCashString)
	var stockValue int
	var usableCash int

	message := quoteRequest(userId, stock, transactionNum)
	stockValueQuoteString := message[0]
	stockValue = stringToCents(stockValueQuoteString)

	//check if user has enough money for a BUY
	if err := sessionGlobal.Query("select userId, usableCash from users where userid='"+userId+"'").Scan(&userId, &usableCash); err != nil {
		panic(err)
	}

	if usableCash < pendingCash {
		log.Printf("Not enough money for user: %s to perform buy", userId)
		return
	}

	//if has enough cash subtract and set aside from db
	usableCash = usableCash - pendingCash
	usableCashString := strconv.Itoa(usableCash)
	pendingCashString = strconv.Itoa(pendingCash)
	//fmt.Println("Available Cash is greater than buy amount")
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(err)
	}
	//fmt.Println("Cash allocated")

	u := uuid.NewV1()
	f := uuid.Formatter(u, uuid.FormatCanonical)
	//fmt.Println(f)

	stockValueString := strconv.FormatInt(int64(stockValue), 10)
	if err := sessionGlobal.Query("INSERT INTO buypendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(err)
	}

	/*
		insert userid and pid into the redis database
		to start decrementing the transaction
		NEED TO HAVE SMOETHING TO CHECK WHEN THE 60 seconds is up
		to return the money back and alert the user
	*/

	//run update function to check if the buy command has expired
	updateStateBuy(f, userId)
}

//checks the state and runs only after a buy or sell to check if the UUID of a transaction is expired or NOT
//this is needed to return the allocated money in the case the transaction automatically expires
//waits for 62 seconds, checks the UUID parameter if it exists in the redis database, and if it doesnt it will revert the buy or sell command
func updateStateBuy(uuid string, userId string) {
	time.Sleep(time.Second * 62)

	//fmt.Println("starting operation 1")
	var pendingCash int
	var usableCash int
	var totalCash int
	var count int

	//check if remaining transaction still exists
	if err := sessionGlobal.Query("select count(*) from buypendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&count); err != nil {
		panic(err)
	}

	// number of pending buy transactions that fit the bill
	if count == 0 {
		fmt.Println("buy transaction doesnt exist")
		return
	}

	//obtain value remaining for expired transaction
	if err := sessionGlobal.Query("select pendingCash from buypendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&pendingCash); err != nil {
		panic(err)
	}

	if pendingCash == 0 {
		return
	}

	//delete pending transaction
	if err := sessionGlobal.Query("delete from buypendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(err)
	}

	//obtain users bank account value
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(err)
	}

	//add back account value to pending cash
	totalCash = pendingCash + usableCash
	totalCashString := strconv.Itoa(totalCash)

	//return total money to user
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(err)
	}

}

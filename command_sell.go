package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/twinj/uuid"
)

func sell(userId string, stock string, sellStockDollarsString string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "SELL", userId, stock, sellStockDollarsString)

	sellStockDollars := stringToCents(sellStockDollarsString)
	var stockValue int
	var usableStocks int
	var stockname string
	var stockamount int
	var usid string
	var hasStock bool

	message := quoteRequest(userId, stock, transactionNum)
	stockValue = stringToCents(message[0])

	//check if user has enough stocks for a SELL
	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
	for iter.Scan(&usid, &stockname, &stockamount) {
		if stockname == stock {
			hasStock = true
			break
		}

	}
	if err := iter.Close(); err != nil {
		panic(err)
	}

	if !hasStock {
		return
	}

	usableStocks = stockamount
	//if not close the session
	if (stockValue * usableStocks) < sellStockDollars {

		return
	}
	if stockValue == 0 {
		return
	}

	sellableStocks := sellStockDollars / stockValue
	usableStocks = usableStocks - sellableStocks
	usableStocksString := strconv.FormatInt(int64(usableStocks), 10)

	pendingCash := sellableStocks * stockValue
	pendingCashString := strconv.FormatInt(int64(pendingCash), 10)
	stockValueString := strconv.FormatInt(int64(stockValue), 10)

	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + usableStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	u := uuid.NewV1()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	if err := sessionGlobal.Query("INSERT INTO sellpendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	updateStateSell(userId, f, usid)
}

func updateStateSell(userId string, uuid string, usid string) {
	//print("In update sell")
	timer1 := time.NewTimer(time.Second * 62)

	<-timer1.C

	var pendingCash int
	var pendingStocks int
	var currentStocks int
	var totalStocks int
	var count int

	//check if remaining transaction still exists
	if err := sessionGlobal.Query("select count(*) from sellpendingtransactions where userid='" + userId + "' and pid=" + uuid).Scan(&count); err != nil {
		panic(err)
	}

	if count == 0 {
		return
	}

	//obtain number of stocks for expired transaction
	if err := sessionGlobal.Query("select pendingcash, stockvalue from sellpendingtransactions where userid='"+userId+"' and pid="+uuid).Scan(&pendingCash, &pendingStocks); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := sessionGlobal.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	//get current users stock amount
	if err := sessionGlobal.Query("select stockamount from userstocks where usid=" + usid).Scan(&currentStocks); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//add back stocks to stocks
	stocks := pendingCash / pendingStocks
	totalStocks = stocks + currentStocks
	totalStocksString := strconv.FormatInt(int64(totalStocks), 10)

	//return total stocks to the userstocks
	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + totalStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(err)
	}

}

package main

import (
//	"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	//"strings"
	"github.com/twinj/uuid"
	"time"
	//"github.com/go-redis/redis"
	//"log"
)

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
	logQuoteEvent(timestamp_quote,"TS1",transactionNum_string,message[0],message[1],userId,message[3],message[4])
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

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	//transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "SELL", userId, stock, sellStockDollarsString)

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


	if err := sessionGlobal.Query("INSERT INTO sellpendingtransactions (pid, userid, pendingCash, stock, stockValue) VALUES (" + f + ", '" + userId + "', " + pendingCashString + ", '" + stock + "' , " + stockValueString + ")").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	
	go updateStateSell(userId, f, usid)

	
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
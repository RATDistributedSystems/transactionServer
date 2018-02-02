package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	//"strings"
	"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
	//"log"
)


//sets the total cash to gain from selling a stock
func setSellAmount(userId string, stock string, pendingCashString string,transactionNum int){


	pendingCashCents := stringToCents(pendingCashString)
	//check if user owns stock
	ownedStockAmount, usid := checkStockOwnership(userId, stock)
	fmt.Println(usid)

	/*
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "SET_SELL_AMOUNT", userId, stock, pendingCashString)
	*/

	if(ownedStockAmount == 0){
		fmt.Println("Cannot Sell a stock you don't own")
		return
	}

	pendingCashString = strconv.FormatInt(int64(pendingCashCents), 10)

	//create trigger to sell a certain amount of the stock
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)


	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)


	//Create new entry for the sell trigger with the sell amount
	if err := sessionGlobal.Query("INSERT INTO sellTriggers (tid, pendingCash, stock, userid) VALUES (" + f + ", " + pendingCashString + ", '" + stock + "', '" + userId + "')").Exec(); err != nil{
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}
}
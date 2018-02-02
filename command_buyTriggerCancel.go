package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	//"strconv"
	//"strings"
	//"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
	//"log"
)


//cancel any buy triggers as well as buy_sell_amounts
func cancelBuyTrigger(userId string, stock string,transactionNum int){
	/*
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "CANCEL_SET_BUY", userId, stock, "")
	*/
	buyExists := checkDependency("CANCEL_SET_BUY",userId,stock)

	if(buyExists == false){
		fmt.Println("cannot CANCEL, no buys pending")
		return
	}

	fmt.Println("cancelling buy trigger")

	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)


	if err := sessionGlobal.Query("DELETE FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}
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


func cancelSellTrigger(userId string, stock string,transactionNum int){

	/*
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "CANCEL_SET_SELL", userId, stock, "")
	*/

	sellExists := checkDependency("CANCEL_SET_SELL",userId,stock)
	if(sellExists == false){
		//fmt.Println("cannot CANCEL, no sells pending")
		return
	}

	//fmt.Println("cancelling sell trigger")


	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)


	if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}
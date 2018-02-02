package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	"strings"
	//"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
	//"log"
)


func cancelBuy(userId string,transactionNum int){
	var pendingCash int
	var usableCash int
	var totalCash int
	var uuid string
	var stock string
	userId = strings.TrimSuffix(userId, "\n")

	/*
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "CANCEL_BUY", userId, stock, "")
	*/

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

	

}
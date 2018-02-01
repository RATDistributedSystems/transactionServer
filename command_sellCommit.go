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
	"time"
	//"github.com/go-redis/redis"
	//"log"
)

func commitSell(userId string,transactionNum int){

	var uuid string
	var pendingCash int
	var usableCash int
	var stock string
	userId = strings.TrimSuffix(userId, "\n")

		timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "COMMIT_SELL", userId, stock, "")

	sellExists := checkDependency("COMMIT_SELL",userId,"none")
	if(sellExists == false){
		fmt.Println("cannot commit, no sells pending")
		return
	}


	if err := sessionGlobal.Query("select pid,stock from sellpendingtransactions where userid='" + userId + "'").Scan(&uuid,&stock); err != nil{
		panic(fmt.Sprintf("problem", err))
	}

	//get pending cash to be added to user account
	if err := sessionGlobal.Query("select pid, pendingcash from sellpendingtransactions where userid='" + userId + "'").Scan(&uuid, &pendingCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if uuid == "" {
		return
	}

	//get current users cash
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	//add available cash to leftover cash
	usableCash = usableCash + pendingCash
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	fmt.Println(usableCashString)

	//re input the new cash value in to the user db
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	//subtract sold stocks from users owned stocks



	//delete the pending transcation
	if err := sessionGlobal.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}



}
package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	//"strings"
	//"github.com/twinj/uuid"
	"time"
	//"github.com/go-redis/redis"
	//"log"
)

func addUser(userId string, usableCashString string,transactionNum int){


	usableCash := stringToCents(usableCashString)

	fmt.Println(usableCash)

	var count int

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logAccountTransactionEvent(timestamp_command, "TS1", "1", "ADD", userId, usableCashString)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "ADD", userId, "", usableCashString)

	if err := sessionGlobal.Query("SELECT count(*) FROM users WHERE userid='" + userId + "'").Scan(&count); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//if the user already exists add money to the account
	if count != 0{
		fmt.Println("adding funds to user")
		addFunds(userId, usableCash)


	//if the user doesnt exist create a new user
	}else{
		fmt.Println("creating new user")
		usableCashString = strconv.FormatInt(int64(usableCash), 10)
		if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('" + userId + "', " + usableCashString + ")").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}	
	}

	//usableCashString = strconv.FormatInt(int64(usableCash), 10)

	//if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('Jones', 351) IF NOT EXISTS").Exec(); err != nil {
	//	panic(fmt.Sprintf("problem creating session", err))
	//}



	//if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('" + userId + "', " + usableCashString + ")").Exec(); err != nil {
	//	panic(fmt.Sprintf("problem creating session", err))
	//}

	
}







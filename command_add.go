package main

import (
	"fmt"
	"strconv"
	"time"
)

func addUser(userId string, usableCashString string, transactionNum int) {

	usableCash := stringToCents(usableCashString)
	var count int

	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	logAccountTransactionEvent(timestamp_command, "TS1", "1", "ADD", userId, usableCashString)

	if err := sessionGlobal.Query("SELECT count(*) FROM users WHERE userid='" + userId + "'").Scan(&count); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//if the user already exists add money to the account
	if count != 0 {
		addFunds(userId, usableCash)

		//if the user doesnt exist create a new user
	} else {
		usableCashString = strconv.FormatInt(int64(usableCash), 10)
		if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('" + userId + "', " + usableCashString + ")").Exec(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
	}
}

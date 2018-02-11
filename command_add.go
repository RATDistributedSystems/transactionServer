package main

import (
	"strconv"
)

func addUser(userID string, usableCashString string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "ADD", userID, "", usableCashString)
	logAccountTransactionEvent("TS1", transactionNum, "ADD", userID, usableCashString)
	usableCash := stringToCents(usableCashString)
	var count int

	if err := sessionGlobal.Query("SELECT count(*) FROM users WHERE userid='" + userID + "'").Scan(&count); err != nil {
		panic(err)
	}

	//if the user already exists add money to the account
	if count != 0 {
		addFunds(userID, usableCash)

		//if the user doesnt exist create a new user
	} else {
		usableCashString = strconv.FormatInt(int64(usableCash), 10)
		if err := sessionGlobal.Query("INSERT INTO users (userid, usableCash) VALUES ('" + userID + "', " + usableCashString + ")").Exec(); err != nil {
			panic(err)
		}
	}
}

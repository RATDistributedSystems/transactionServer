package main

import (
	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

func addUser(userID string, usableCashString string, transactionNum int) {
	//logUserEvent("TS1", transactionNum, "ADD", userID, "", usableCashString)
	logAccountTransactionEvent(serverName, transactionNum, "ADD", userID, usableCashString)

	addFundAmount := stringToCents(usableCashString)
	if ratdatabase.UserExists(userID) {
		balance := ratdatabase.GetUserBalance(userID)
		newBalance := balance + addFundAmount
		ratdatabase.UpdateUserBalance(userID, newBalance)
	} else { // user doesn't exist yet
		ratdatabase.CreateUser(userID, addFundAmount)
	}
}

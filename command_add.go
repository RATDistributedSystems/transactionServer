package main

import (
	"fmt"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandAdd struct {
	username string
	amount   string
}

func (c commandAdd) process(transaction int) {
	logUserEvent(serverName, transaction, "ADD", c.username, "", c.amount)
	addUser(c.username, c.amount, transaction)
}

func addUser(userID string, usableCashString string, transactionNum int) {
	fmt.Printf("[%d] ADD,%s,%s", transactionNum, userID, usableCashString)
	logAccountTransactionEvent(serverName, transactionNum, "ADD", userID, usableCashString)
	fmt.Println(usableCashString)
	addFundAmount := stringToCents(usableCashString)
	if ratdatabase.UserExists(userID) {
		balance := ratdatabase.GetUserBalance(userID)
		newBalance := balance + addFundAmount
		ratdatabase.UpdateUserBalance(userID, newBalance)
	} else { // user doesn't exist yet
		ratdatabase.CreateUser(userID, addFundAmount)
	}
}

package main

import (
	"fmt"
	"time"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandAdd struct {
	username string
	amount   string
}

func (c commandAdd) process(transaction int) string {
	logUserEvent(serverName, transaction, "ADD", c.username, "", c.amount)
	return addUser(c.username, c.amount, transaction)
}

func addUser(userID string, usableCashString string, transactionNum int) string {

	start := time.Now()

	fmt.Printf("[%d] ADD,%s,%s", transactionNum, userID, usableCashString)
	logAccountTransactionEvent(serverName, transactionNum, "ADD", userID, usableCashString)
	addFundAmount := stringToCents(usableCashString)
	start2 := time.Now()
	if ratdatabase.UserExists(userID) {
		elapsed2 := time.Since(start2)
		appendToText("addCQL.txt", elapsed2.String())
		balance := ratdatabase.GetUserBalance(userID)
		newBalance := balance + addFundAmount
		ratdatabase.UpdateUserBalance(userID, newBalance)
		elapsed := time.Since(start)
		appendToText("add.txt", elapsed.String())
		return fmt.Sprintf("Added %s to %s's account", usableCashString, userID)
	}

	// user doesn't exist yet
	start1 := time.Now()
	ratdatabase.CreateUser(userID, addFundAmount)
	elapsed1 := time.Since(start1)
	appendToText("addCQL.txt", elapsed1.String())

	elapsed := time.Since(start)
	appendToText("add.txt", elapsed.String())

	return fmt.Sprintf("Added new user %s with starting balance %s", userID, usableCashString)
}

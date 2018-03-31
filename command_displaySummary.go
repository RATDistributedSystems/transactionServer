package main

import (
	"fmt"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandDisplaySummary struct {
	username string
}

func (c commandDisplaySummary) process(transaction int) {
	logUserEvent(serverName, transaction, "DISPLAY_SUMMARY", c.username, "", "")
	displaySummary(c.username, transaction)
}

func displaySummary(userId string, transactionNum int) {
	//get user cash
	fmt.Println(userId + "summary:")
	currentBalance := ratdatabase.GetUserBalance(userId)
	fmt.Printf("Balance: %n", currentBalance)
	//get list of user stocks and their quantity
	fmt.Println("Stock summary: ")
	rs := ratdatabase.GetStockAndAmountOwned(userId)
	for _, r := range rs {
		fmt.Println(r["stock"])
		fmt.Println(r["stockamount"])
	}
	//get list of buy triggers

	//get list of sell triggers

}

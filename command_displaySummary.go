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
	fmt.Println("Buy Triggers set: ")
	bt := ratdatabase.GetBuyTriggers(userId)
	for _, b := range bt {
		fmt.Println(b["stock"])
		fmt.Println(b["stockamount"])
	}

	//get list of sell triggers
	fmt.Println("Sell Triggers set: ")
	st := ratdatabase.GetSellTriggers(userId)
	for _, s := range st{
		fmt.Println(s["stock"])
		fmt.Println(s["stockamount"])
	}

}

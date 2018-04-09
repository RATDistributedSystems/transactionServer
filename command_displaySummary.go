package main

import (
	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandDisplaySummary struct {
	username string
}

func (c commandDisplaySummary) process(transaction int) string {
	logUserEvent(serverName, transaction, "DISPLAY_SUMMARY", c.username, "", "")
	return displaySummary(c.username, transaction)
}

func displaySummary(userId string, transactionNum int) string {
	if !ratdatabase.UserExists(userId) {
		return "No user exists"
	}

	account := userAcount{}
	account.name = userId
	account.balance = ratdatabase.GetUserBalance(userId)

	//get list of user stocks and their quantity
	rs := ratdatabase.GetStockAndAmountOwned(userId)
	for _, r := range rs {
		name := r["stock"].(string)
		amount := r["stockamount"].(int)
		stock := stockDetails{name, amount}
		account.stocks = append(account.stocks, stock)
	}

	//get list of buy triggers
	bt := ratdatabase.GetBuyTriggers(userId)
	for _, b := range bt {
		name := b["stock"].(string)
		amount := b["stockamount"].(int)
		trigger := stockTriggerDetails{name, amount, 0}
		account.buytriggers = append(account.buytriggers, trigger)
	}

	//get list of sell triggers
	st := ratdatabase.GetSellTriggers(userId)
	for _, s := range st {
		name := s["stock"].(string)
		amount := s["stockamount"].(int)
		trigger := stockTriggerDetails{name, amount, 0}
		account.selltriggers = append(account.selltriggers, trigger)
	}

	return account.String()
}

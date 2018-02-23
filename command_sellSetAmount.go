package main

import (
	"fmt"
	"strconv"

	"github.com/twinj/uuid"
)

//sets the total cash to gain from selling a stock
func setSellAmount(userId string, stock string, pendingCashString string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "SET_SELL_AMOUNT", userId, stock, pendingCashString)
	pendingCashCents := stringToCents(pendingCashString)
	//check if user owns stock
	ownedStockAmount, _ := checkStockOwnership(userId, stock)

	if ownedStockAmount == 0 {
		fmt.Println("Cannot Sell a stock you don't own")
		return
	}

	//check amount of stock that can be sold
	//divide the pendingCashString by the current quote price to determine the amount of total sellable stocks

	//check current stock price
	currentStockPrice := quoteRequest(userId, stock, transactionNum)

	//make sure amount wanting to buy isnt too high
	if currentStockPrice > pendingCashCents {
		fmt.Println("Current stock price is greater than amount wanting to be sold")
		return
	}

	sellStockName := stock

	var usid string
	var ownedstockname string
	var stockamount int

	//check how much of the stock the user has
	iter := sessionGlobal.Query("SELECT usid, stock, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
	for iter.Scan(&usid, &ownedstockname, &stockamount) {
		if ownedstockname == sellStockName {
			break
		}
		//fmt.Println("STOCKS: ", stock, stockamount)
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//calculate amount of stocks can be bought
	sellStockAmount := pendingCashCents / currentStockPrice
	//subtract stocks to allocate from owned stocks
	sellStockAmountString := strconv.FormatInt(int64(sellStockAmount), 10)
	stockLeftOver := stockamount - sellStockAmount
	stockLeftOverString := strconv.FormatInt(int64(stockLeftOver), 10)

	//allocate these stocks
	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount=" + stockLeftOverString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//create trigger to sell a certain amount of the stock
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	//Create new entry for the sell trigger with the sell amount
	if err := sessionGlobal.Query("INSERT INTO sellTriggers (tid, pendingStocks, stock, userid) VALUES (" + f + ",'" + sellStockAmountString + "', '" + stock + "', '" + userId + "')").Exec(); err != nil {
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}
}

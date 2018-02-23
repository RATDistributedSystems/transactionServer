package main

import (
	"fmt"
	"strconv"
)

func cancelSellTrigger(userId string, stock string, transactionNum int) {
	//logUserEvent("TS1", transactionNum, "CANCEL_SET_SELL", userId, stock, "")
	sellExists := checkDependency("CANCEL_SET_SELL", userId, stock)
	if sellExists == false {
		//fmt.Println("cannot CANCEL, no sells pending")
		return
	}

	sellStockName := stock

	//get user stocks
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
		panic(err)
	}

	//get stocks allocated to sell
	var pendingStocks int
	if err := sessionGlobal.Query("SELECT pendingStocks FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "' ").Scan(&pendingStocks); err != nil {
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}

	//re add allocated trigger stocks and stocks the user has
	stockTotal := stockamount + pendingStocks
	stockTotalString := strconv.FormatInt(int64(stockTotal), 10)
	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount=" + stockTotalString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(err)
	}
}

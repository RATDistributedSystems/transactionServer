package main

import (
	"fmt"
)

//cancel any buy triggers as well as buy_sell_amounts
func cancelBuyTrigger(userId string, stock string, transactionNum int) {
	//logUserEvent("TS1", transactionNum, "CANCEL_SET_BUY", userId, stock, "")
	buyExists := checkDependency("CANCEL_SET_BUY", userId, stock)

	if buyExists == false {
		fmt.Println("cannot CANCEL, no buys pending")
		return
	}

	if err := sessionGlobal.Query("DELETE FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}

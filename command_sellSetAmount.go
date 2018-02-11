package main

import (
	"fmt"
	"strconv"

	"github.com/twinj/uuid"
)

//sets the total cash to gain from selling a stock
func setSellAmount(userId string, stock string, pendingCashString string, transactionNum int) {

	pendingCashCents := stringToCents(pendingCashString)
	//check if user owns stock
	ownedStockAmount, _ := checkStockOwnership(userId, stock)

	if ownedStockAmount == 0 {
		fmt.Println("Cannot Sell a stock you don't own")
		return
	}

	pendingCashString = strconv.FormatInt(int64(pendingCashCents), 10)

	//create trigger to sell a certain amount of the stock
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	//Create new entry for the sell trigger with the sell amount
	if err := sessionGlobal.Query("INSERT INTO sellTriggers (tid, pendingCash, stock, userid) VALUES (" + f + ", " + pendingCashString + ", '" + stock + "', '" + userId + "')").Exec(); err != nil {
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}
}

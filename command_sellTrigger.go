package main

import (
	"fmt"
	"strconv"
	"time"
)

func setSellTrigger(userId string, stock string, stockSellPrice string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "SET_SELL_TRIGGER", userId, stock, stockSellPrice)
	stockSellPriceCents := stringToCents(stockSellPrice)
	stockSellPriceCentsString := strconv.FormatInt(int64(stockSellPriceCents), 10)

	//check if setSellAmount is set for this particular stock
	var count int
	if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "' ").Scan(&count); err != nil {
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}
	//if set sell amount isnt set return
	if count == 0 {
		return
	}

	//update database entry with trigger value
	if err := sessionGlobal.Query("UPDATE sellTriggers SET triggerValue=" + stockSellPriceCentsString + " WHERE userid='" + userId + "' AND stock='" + stock + "' ").Exec(); err != nil {
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}

	go checkSellTrigger(userId, stock, stockSellPriceCents, transactionNum)
}

func checkSellTrigger(userId string, stock string, stockSellPriceCents int, transactionNum int) {

	operation := false

	for {
		//check the quote server every 5 seconds
		timer1 := time.NewTimer(time.Second * 1)
		<-timer1.C

		//if the trigger doesnt exist exit

		exists := checkTriggerExists(userId, stock, operation)
		if exists == false {
			return
		}

		//retrieve current stock price
		currentStockPrice := quoteRequest(userId, stock, transactionNum)

		if currentStockPrice > stockSellPriceCents {

			//sell the allocated stocks

			//get stocks allocated to sell
			var pendingStocks int
			if err := sessionGlobal.Query("SELECT pendingStocks FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "' ").Scan(&pendingStocks); err != nil {
				//panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
				return
			}

			sellProfits := pendingStocks * currentStockPrice

			//delete pending transaction
			if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "' ").Exec(); err != nil {
				//panic(err)
				return
			}

			//add profits from selling stock to account
			fmt.Println("Sell Trigger Sucessful, profits added to account")
			addFunds(userId, sellProfits)
			return
		}

	}

}

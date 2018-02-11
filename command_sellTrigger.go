package main

import (
	"fmt"
	"strconv"
	"time"
)

func setSellTrigger(userId string, stock string, stockSellPrice string, transactionNum int) {

	stockSellPriceCents := stringToCents(stockSellPrice)
	stockSellPriceCentsString := strconv.FormatInt(int64(stockSellPriceCents), 10)

	//check if set sell amount is set for this particular stock
	var count int
	if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "' ").Scan(&count); err != nil {
		panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
	}

	//if set sell amount isnt set return
	if count == 0 {
		//fmt.Println("No set Sell amount placed")
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
		timer1 := time.NewTimer(time.Second * 5)
		<-timer1.C

		//if the trigger doesnt exist exit

		exists := checkTriggerExists(userId, stock, operation)
		if exists == false {
			return
		}

		//retrieve current stock price
		message := quoteRequest(userId, stock, transactionNum)
		currentStockPrice := stringToCents(message[0])

		if currentStockPrice > stockSellPriceCents {

			//Check how many stocks the user can sell

			var pendingCash int

			if err := sessionGlobal.Query("SELECT pendingCash FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "' ").Scan(&pendingCash); err != nil {
				panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
			}

			//calculate amount of stocks can be sold
			sellAbleStocksMax := pendingCash / currentStockPrice

			//check how many stocks the user owns
			ownedStocks, usid := checkStockOwnership(userId, stock)

			//check how many stocks the user can sell
			var sellAbleStocks int
			var remainingStocks int

			//check if user has more owned stocks than able to sell
			if sellAbleStocksMax < ownedStocks {
				sellAbleStocks = sellAbleStocksMax
				remainingStocks = ownedStocks - sellAbleStocks

				//calculate money gained from stocks selling
				sellAbleStockPrice := pendingCash - (sellAbleStocksMax * currentStockPrice)

				remainingStocksString := strconv.FormatInt(int64(remainingStocks), 10)

				//update userStock database with new about of stock
				if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + remainingStocksString + " WHERE usid=" + usid).Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//increase money in userAccount
				addFunds(userId, sellAbleStockPrice)

				//delete trigger

				if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
			} else {
				//Case where user does not own enough stocks to sell the maximum amount
				//user must sell the most it can
				sellAbleStocks = ownedStocks

				//ownedStocksString, err := strconv.Atoi("ownedStocks")

				sellAbleStockPrice := pendingCash - (sellAbleStocks * currentStockPrice)
				remainingCash := pendingCash - sellAbleStockPrice
				sellAbleStockPrice = sellAbleStockPrice + remainingCash

				addFunds(userId, sellAbleStockPrice)

				if err := sessionGlobal.Query("DELETE FROM userstocks WHERE usid=" + usid).Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}

				//delete trigger
				if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
			}
		}

	}

}

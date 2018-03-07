package main

import (
	"fmt"
	"log"
	"time"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

func setSellAmount(userID string, stock string, pendingCashString string, transactionNum int) {
	//logUserEvent("TS1", transactionNum, "SET_SELL_AMOUNT", userId, stock, pendingCashString)
	stockAmountToSell := stringToCents(pendingCashString)
	//check if user owns stock

	uuid, stockAmountOwned, ownsStock := ratdatabase.GetStockAmountOwned(userID, stock)

	if !ownsStock || stockAmountOwned == 0 {
		msg := "[%d] Not enough of stock %s (%d) for SellSetAmount %d for %s"
		log.Printf(msg, transactionNum, stock, stockAmountOwned, stockAmountToSell, userID)
		return
	}

	currentStockPrice := quoteRequest(userID, stock, transactionNum)

	if currentStockPrice > stockAmountToSell {
		msg := "[%d] Current stock price (%d) is greater than amount wanting to be sold (%d)"
		fmt.Printf(msg, transactionNum, currentStockPrice, stockAmountToSell)
		return
	}

	sellStockAmount := stockAmountToSell / currentStockPrice

	if stockAmountOwned < sellStockAmount {
		msg := "[%d] User %s does not own enough %s stock (%d) to sell %d amount"
		fmt.Printf(msg, transactionNum, userID, stock, stockAmountOwned, sellStockAmount)
		return
	}

	// Insert stock amount and return stock if it was already made
	oldSellStockAmount := ratdatabase.InsertSetSellTrigger(userID, stock, sellStockAmount)
	remainingStockAmount := stockAmountOwned + oldSellStockAmount - sellStockAmount
	ratdatabase.UpdateUserStockByUUID(uuid, stock, remainingStockAmount)
}

func setSellTrigger(userID string, stock string, stockSellPrice string, transactionNum int) {
	//logUserEvent("TS1", transactionNum, "SET_SELL_TRIGGER", userId, stock, stockSellPrice)
	stockValueCents := stringToCents(stockSellPrice)
	stockAmountSet := ratdatabase.UpdateSellTriggerPrice(userID, stock, stockValueCents)

	if !stockAmountSet {
		msg := "[%d] User %s hasn't set stock amount for stock %s"
		log.Printf(msg, transactionNum, userID, stock)
		return
	}

	checkSellTrigger(userID, stock, stockValueCents, transactionNum)
}

func cancelSellTrigger(userID string, stock string, transactionNum int) {
	//logUserEvent("TS1", transactionNum, "CANCEL_SET_SELL", userId, stock, "")
	returnAmount := ratdatabase.CancelSellTrigger(userID, stock)

	if returnAmount == 0 {
		return
	}

	uuid, oldStockAmount, _ := ratdatabase.GetStockAmountOwned(userID, stock)
	newStockAmount := oldStockAmount + returnAmount
	ratdatabase.UpdateUserStockByUUID(uuid, stock, newStockAmount)
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
			if err := sessionGlobal.Query("SELECT stockAmount FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "' ").Scan(&pendingStocks); err != nil {
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

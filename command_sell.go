package main

import (
	"fmt"
	"log"
	"time"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandSell struct {
	username string
	amount   string
	stock    string
}

func (c commandSell) process(transaction int) {
	logUserEvent(serverName, transaction, "SELL", c.username, c.stock, c.amount)
	sell(c.username, c.stock, c.amount, transaction)
}

type commandCanceSell struct {
	username string
}

func (c commandCanceSell) process(transaction int) {
	logUserEvent(serverName, transaction, "CANCEL_SELL", c.username, "", "")
	cancelSell(c.username, transaction)
}

type commandCommitSell struct {
	username string
}

func (c commandCommitSell) process(transaction int) {
	logUserEvent(serverName, transaction, "COMMIT_SELL", c.username, "", "")
	commitSell(c.username, transaction)
}

func sell(userId string, stock string, sellStockDollarsString string, transactionNum int) {
	sellStockValue := stringToCents(sellStockDollarsString)
	//stockValue := quoteRequest(userId, stock, transactionNum)
	stockValue := getQuote(userId, stock, transactionNum)

	// unlikely but has happened before
	if stockValue == 0 {
		log.Printf("[%d] Stock '%s' price is 0. Cannot buy", transactionNum, stock)
		logErrorEvent(serverName, transactionNum, "SELL", userId, stock, sellStockDollarsString, "Stock price is 0. Cannot buy")
		return
	}

	stockToSell := sellStockValue / stockValue
	if stockToSell == 0 {
		stockToSell = 1
	}
	potentialProfit := stockToSell * stockValue

	fmt.Printf("User Req Sell: %d / Stock Price: %d == %d\n", sellStockValue, stockValue, stockToSell)

	_, stockAmount, ownsStock := ratdatabase.GetStockAmountOwned(userId, stock)
	if !ownsStock || stockToSell > stockAmount {
		log.Printf("[%d] %s doesn't have enough stock %s@%.2f to sell. Have: %d, Need: %d", transactionNum, userId, stock, float64(stockValue/100), stockAmount, stockToSell)
		logErrorEvent(serverName, transactionNum, "SELL", userId, stock, sellStockDollarsString, "Not enough stock to sell.")
		return
	}

	// Remove stock from userstocks table and move to pendingselltransaction table
	newStockAmount := stockAmount - stockToSell
	ratdatabase.UpdateUserStockByUserAndStock(userId, stock, newStockAmount)
	//ratdatabase.UpdateUserStockByUUID(userUUID, stock, newStockAmount)
	transactionUUID := ratdatabase.InsertPendingSellTransaction(userId, stock, potentialProfit, stockValue)
	log.Printf("[%d] User %s sell transaction for %d %s@%.2f pending", transactionNum, userId, stockToSell, stock, float64(stockValue/100))

	time.Sleep(time.Second * 62)

	// If transaction isn't alive, nothing to do
	if !ratdatabase.SellTransactionAlive(userId, transactionUUID) {
		return
	}
	log.Printf("[%d] Cancelling '%s' request to sell %.2f of stock %s\n", transactionNum, userId, float64(sellStockValue/100), stock)

	// delete pending transaction
	ratdatabase.DeletePendingSellTransaction(userId, transactionUUID)

	// Add stocks back to account
	_, currentStockAmount, _ := ratdatabase.GetStockAmountOwned(userId, stock)

	latestStockAmount := currentStockAmount + stockToSell
	//ratdatabase.UpdateUserStockByUUID(userUUID, stock, latestStockAmount)
	ratdatabase.UpdateUserStockByUserAndStock(userId, stock, latestStockAmount)

}

func cancelSell(userID string, transactionNum int) {
	transactionUUID, profits, stockName, stockPrice, exists := ratdatabase.GetLastPendingSellTransaction(userID)

	if !exists {
		log.Printf("[%d] No pending sell transaction to cancel", transactionNum)
		logErrorEvent(serverName, transactionNum, "CANCEL_SELL", userID, "", "", "No pending sell transaction to cancel.")
		return
	}

	// Return the stock
	stockToReturn := profits / stockPrice
	_, stockAmount, owmsStock := ratdatabase.GetStockAmountOwned(userID, stockName)

	if owmsStock {
		newStockAmount := stockAmount + stockToReturn
		ratdatabase.UpdateUserStockByUserAndStock(userID, stockName, newStockAmount)
		//ratdatabase.UpdateUserStockByUUID(stockUUID, stockName, newStockAmount)
	} else {
		ratdatabase.AddStockToPortfolio(userID, stockName, stockToReturn)
	}

	//delete pending transaction
	ratdatabase.DeletePendingSellTransaction(userID, transactionUUID)
}

func commitSell(userId string, transactionNum int) {
	transactionUUID, profits, stockName, stockPrice, exists := ratdatabase.GetLastPendingSellTransaction(userId)

	if !exists {
		log.Printf("[%d] No pending sell transaction to commit", transactionNum)
		logErrorEvent(serverName, transactionNum, "COMMIT_SELL", userId, "", "", "No pending sell transaction to commit.")
		return
	}

	log.Printf("[%d] Commiting sale of %s@%.2f for %.2f for User %s", transactionNum, stockName, float64(stockPrice/100), float64(profits/100), userId)

	currentBalance := ratdatabase.GetUserBalance(userId)
	newBalance := currentBalance + profits
	ratdatabase.UpdateUserBalance(userId, newBalance)

	//delete the pending transcation
	ratdatabase.DeletePendingSellTransaction(userId, transactionUUID)
}

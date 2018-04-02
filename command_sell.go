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

func (c commandSell) process(transaction int) string {
	logUserEvent(serverName, transaction, "SELL", c.username, c.stock, c.amount)
	return sell(c.username, c.stock, c.amount, transaction)
}

type commandCanceSell struct {
	username string
}

func (c commandCanceSell) process(transaction int) string {
	logUserEvent(serverName, transaction, "CANCEL_SELL", c.username, "", "")
	return cancelSell(c.username, transaction)
}

type commandCommitSell struct {
	username string
}

func (c commandCommitSell) process(transaction int) string {
	logUserEvent(serverName, transaction, "COMMIT_SELL", c.username, "", "")
	return commitSell(c.username, transaction)
}

func sell(userId string, stock string, sellStockDollarsString string, transactionNum int) string {
	sellStockValue := stringToCents(sellStockDollarsString)
	//stockValue := quoteRequest(userId, stock, transactionNum)
	stockValue := getQuote(userId, stock, transactionNum)

	// unlikely but has happened before
	if stockValue == 0 {
		m := fmt.Sprintf("[%d] Stock '%s' price is 0. Cannot buy", transactionNum, stock)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "SELL", userId, stock, sellStockDollarsString, m)
		return m
	}

	stockToSell := sellStockValue / stockValue
	if stockToSell == 0 {
		stockToSell = 1
	}
	potentialProfit := stockToSell * stockValue

	//fmt.Printf("User Req Sell: %d / Stock Price: %d == %d\n", sellStockValue, stockValue, stockToSell)

	_, stockAmount, ownsStock := ratdatabase.GetStockAmountOwned(userId, stock)
	if !ownsStock || stockToSell > stockAmount {
		m := fmt.Sprintf("[%d] %s does not have enough stock %s@%.2f to sell. Have: %d, Need: %d", transactionNum, userId, stock, float64(stockValue/100), stockAmount, stockToSell)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "SELL", userId, stock, sellStockDollarsString, m)
		return m
	}

	// Remove stock from userstocks table and move to pendingselltransaction table
	newStockAmount := stockAmount - stockToSell
	ratdatabase.UpdateUserStockByUserAndStock(userId, stock, newStockAmount)
	//ratdatabase.UpdateUserStockByUUID(userUUID, stock, newStockAmount)
	transactionUUID := ratdatabase.InsertPendingSellTransaction(userId, stock, potentialProfit, stockValue)
	log.Printf("[%d] User %s sell transaction for %d %s@%.2f pending", transactionNum, userId, stockToSell, stock, float64(stockValue/100))
	go checkSell(userId, transactionUUID, stock, stockToSell, sellStockValue, transactionNum)
	return fmt.Sprintf("Buy request for %s@%.2f pending commit/cancel", stock, float64(sellStockValue/100))
}

func checkSell(userID string, transactionUUID string, stock string, stockToSell int, sellStockValue int, transactionNum int) {
	time.Sleep(time.Second * 62)

	// If transaction isn't alive, nothing to do
	if !ratdatabase.SellTransactionAlive(userID, transactionUUID) {
		return
	}
	log.Printf("[%d] Cancelling '%s' request to sell %.2f of stock %s\n", transactionNum, userID, float64(sellStockValue/100), stock)

	// delete pending transaction
	ratdatabase.DeletePendingSellTransaction(userID, transactionUUID)

	// Add stocks back to account
	_, currentStockAmount, _ := ratdatabase.GetStockAmountOwned(userID, stock)

	latestStockAmount := currentStockAmount + stockToSell
	//ratdatabase.UpdateUserStockByUUID(userUUID, stock, latestStockAmount)
	ratdatabase.UpdateUserStockByUserAndStock(userID, stock, latestStockAmount)

}

func cancelSell(userID string, transactionNum int) string {
	transactionUUID, profits, stockName, stockPrice, exists := ratdatabase.GetLastPendingSellTransaction(userID)

	if !exists {
		m := fmt.Sprintf("[%d] No pending sell transaction to cancel", transactionNum)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "CANCEL_SELL", userID, "", "", m)
		return m
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
	return fmt.Sprintf("Cancelled Sale of %s", stockName)
}

func commitSell(userId string, transactionNum int) string {
	transactionUUID, profits, stockName, stockPrice, exists := ratdatabase.GetLastPendingSellTransaction(userId)

	if !exists {
		m := fmt.Sprintf("[%d] No pending sell transaction to commit", transactionNum)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "COMMIT_SELL", userId, "", "", "No pending sell transaction to commit.")
		return m
	}

	log.Printf("[%d] Commiting sale of %s@%.2f for %.2f for User %s", transactionNum, stockName, float64(stockPrice/100), float64(profits/100), userId)

	currentBalance := ratdatabase.GetUserBalance(userId)
	newBalance := currentBalance + profits
	ratdatabase.UpdateUserBalance(userId, newBalance)

	//delete the pending transcation
	ratdatabase.DeletePendingSellTransaction(userId, transactionUUID)
	return fmt.Sprintf("Commited sale of %s", stockName)
}

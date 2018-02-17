package main

import (
	"log"
	"time"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

func sell(userId string, stock string, sellStockDollarsString string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "SELL", userId, stock, sellStockDollarsString)

	sellStockValue := stringToCents(sellStockDollarsString)
	stockValue := quoteRequest(userId, stock, transactionNum)

	// unlikely but has happened before
	if stockValue == 0 {
		log.Printf("[%d] Stock '%s' price is 0. Cannot buy", transactionNum, stock)
		return
	}
	stockToSell := sellStockValue / stockValue
	potentialProfit := stockToSell * stockValue

	userUUID, stockAmount, ownsStock := ratdatabase.GetStockAmountOwned(userId, stock)
	if !ownsStock || stockToSell > stockAmount {
		log.Printf("[%d] %s doesn't have enough stock '%s' to sell. Not proceeding with sell", transactionNum, userId, stock)
		return
	}

	// Remove stock from userstocks table and move to pendingselltransaction table
	newStockAmount := stockAmount - stockToSell
	ratdatabase.UpdateUserStockByUUID(userUUID, stock, newStockAmount)
	transactionUUID := ratdatabase.InsertPendingSellTransaction(userId, stock, potentialProfit, stockValue)
	log.Printf("[%d] User %s sell transaction for %d %s@%d pending", transactionNum, userId, stockAmount, stock, stockValue)

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
	ratdatabase.UpdateUserStockByUUID(userUUID, stock, latestStockAmount)

}

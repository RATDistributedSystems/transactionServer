package main

import (
	"log"
	"time"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

func buy(userId string, stock string, pendingCashString string, transactionNum int) {

	logUserEvent("TS1", transactionNum, "BUY", userId, stock, pendingCashString)

	pendingTransactionCash := stringToCents(pendingCashString)
	message := quoteRequest(userId, stock, transactionNum)
	stockValue := stringToCents(message[0])
	stockAmount := pendingTransactionCash / stockValue
	currentBalance := ratdatabase.GetUserBalance(userId)

	if currentBalance < pendingTransactionCash || stockAmount == 0 {
		log.Printf("[%d] Not enough money for %s to perform buy", transactionNum, userId)
		return
	}

	//if has enough cash subtract and set aside from db
	newBalance := currentBalance - pendingTransactionCash
	ratdatabase.UpdateUserBalance(userId, newBalance)
	uuid := ratdatabase.InsertPendingBuyTransaction(userId, pendingTransactionCash, stock, stockValue)

	//waits for 62 seconds and checks if the transaction is still there. Remove if it is
	time.Sleep(time.Second * 62)

	// If Transaction isn't alive, do nothing
	if !ratdatabase.BuyTransactionAlive(userId, uuid) {
		return
	}

	log.Printf("[%d] Cancelling '%s' request to buy %.2f of stock %s\n", transactionNum, userId, float64(pendingTransactionCash/100), stock)
	ratdatabase.DeletePendingBuyTransaction(userId, uuid)

	// Returns users cash being held
	newerBalance := ratdatabase.GetUserBalance(userId)
	newererBalance := pendingTransactionCash + newerBalance
	ratdatabase.UpdateUserBalance(userId, newererBalance)
}

func cancelBuy(userId string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "CANCEL_BUY", userId, "", "")

	uuid, holdingCash, stockName, _, exists := ratdatabase.GetLastPendingBuyTransaction(userId)

	if !exists {
		log.Printf("[%d] Cannot cancel buy. No buys pending", transactionNum)
		return
	}

	activeBalance := ratdatabase.GetUserBalance(userId)
	newBalance := activeBalance + holdingCash
	ratdatabase.UpdateUserBalance(userId, newBalance)

	//delete pending transaction
	ratdatabase.DeletePendingBuyTransaction(userId, uuid)
	log.Printf("[%d] Buy for stock:%s cancelled by user.", transactionNum, stockName)
}

func commitBuy(userID string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "COMMIT_BUY", userID, "", "")

	uuid, holdingCash, stockName, stockPrice, exists := ratdatabase.GetLastPendingBuyTransaction(userID)

	if !exists {
		log.Printf("[%d] Cannot commit buy for %s. No buy pending", transactionNum, userID)
		return
	}

	stockBought := holdingCash / stockPrice
	surplusCash := holdingCash - (stockBought * stockPrice)

	stockUUID, stockAmount, hasStock := ratdatabase.GetStockAmountOwned(userID, stockName)

	if hasStock {
		newStockAmount := stockAmount + stockBought
		ratdatabase.UpdateUserStockByUUID(stockUUID, stockName, newStockAmount)
	} else {
		ratdatabase.AddStockToPortfolio(userID, stockName, stockAmount)
	}

	if surplusCash != 0 {
		currentBalance := ratdatabase.GetUserBalance(userID)
		newBalance := currentBalance + surplusCash
		ratdatabase.UpdateUserBalance(userID, newBalance)
	}

	ratdatabase.DeletePendingBuyTransaction(userID, uuid)

}

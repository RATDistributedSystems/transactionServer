package main

import (
	"log"
	"time"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandBuy struct {
	username string
	amount   string
	stock    string
}

func (c commandBuy) process(transaction int) {
	logUserEvent(serverName, transaction, "BUY", c.username, c.stock, c.amount)
	buy(c.username, c.stock, c.amount, transaction)
}

type commandCancelBuy struct {
	username string
}

func (c commandCancelBuy) process(transaction int) {
	logUserEvent(serverName, transaction, "CANCEL_BUY", c.username, "", "")
	cancelBuy(c.username, transaction)
}

type commandCommitBuy struct {
	username string
}

func (c commandCommitBuy) process(transaction int) {
	logUserEvent(serverName, transaction, "COMMIT_BUY", c.username, "", "")
	commitBuy(c.username, transaction)
}

func buy(userID string, stock string, pendingCashString string, transactionNum int) {
	pendingTransactionCash := stringToCents(pendingCashString)
	stockValue := getQuote(userID, stock, transactionNum)
	if stockValue <= 0 {
		return
	}
	stockAmount := pendingTransactionCash / stockValue
	currentBalance := ratdatabase.GetUserBalance(userID)

	if currentBalance < pendingTransactionCash {
		log.Printf("[%d] Not enough money for %s to perform buy", transactionNum, userID)
		return
	}

	if stockAmount == 0 {
		log.Printf("[%d] %s stock price(%d) higher than amount to purchase(%d)", transactionNum, stock, stockValue, pendingTransactionCash)
		return
	}

	//if has enough cash subtract and set aside from db
	newBalance := currentBalance - pendingTransactionCash
	ratdatabase.UpdateUserBalance(userID, newBalance)
	uuid := ratdatabase.InsertPendingBuyTransaction(userID, pendingTransactionCash, stock, stockValue)
	log.Printf("[%d] User %s buy transaction for %d %s@%.2f pending", transactionNum, userID, stockAmount, stock, float64(stockValue))

	//waits for 62 seconds and checks if the transaction is still there. Remove if it is
	time.Sleep(time.Second * 62)

	// If Transaction isn't alive, do nothing
	if !ratdatabase.BuyTransactionAlive(userID, uuid) {
		return
	}

	log.Printf("[%d] Cancelling '%s' request to buy %.2f of stock %s\n", transactionNum, userID, float64(pendingTransactionCash/100), stock)
	ratdatabase.DeletePendingBuyTransaction(userID, uuid)

	// Returns users cash being held
	newerBalance := ratdatabase.GetUserBalance(userID)
	newererBalance := pendingTransactionCash + newerBalance
	ratdatabase.UpdateUserBalance(userID, newererBalance)
}

func cancelBuy(userID string, transactionNum int) {
	uuid, holdingCash, stockName, _, exists := ratdatabase.GetLastPendingBuyTransaction(userID)

	if !exists {
		log.Printf("[%d] Cannot cancel buy. No buys pending", transactionNum)
		return
	}

	activeBalance := ratdatabase.GetUserBalance(userID)
	newBalance := activeBalance + holdingCash
	ratdatabase.UpdateUserBalance(userID, newBalance)

	//delete pending transaction
	ratdatabase.DeletePendingBuyTransaction(userID, uuid)
	log.Printf("[%d] Buy for stock:%s cancelled by user.", transactionNum, stockName)
}

func commitBuy(userID string, transactionNum int) {
	uuid, holdingCash, stockName, stockPrice, exists := ratdatabase.GetLastPendingBuyTransaction(userID)

	if !exists {
		log.Printf("[%d] Cannot commit buy for %s. No buy pending", transactionNum, userID)
		return
	}

	stockBought := holdingCash / stockPrice
	surplusCash := holdingCash - (stockBought * stockPrice)

	_, stockAmount, hasStock := ratdatabase.GetStockAmountOwned(userID, stockName)

	if hasStock {
		newStockAmount := stockAmount + stockBought
		ratdatabase.UpdateUserStockByUserAndStock(userID, stockName, newStockAmount)
	} else {
		ratdatabase.AddStockToPortfolio(userID, stockName, stockBought)
	}
	log.Printf("[%d] User %s now has %d more of stock %s", transactionNum, userID, stockBought, stockName)

	if surplusCash != 0 {
		currentBalance := ratdatabase.GetUserBalance(userID)
		newBalance := currentBalance + surplusCash
		ratdatabase.UpdateUserBalance(userID, newBalance)
	}

	ratdatabase.DeletePendingBuyTransaction(userID, uuid)

}

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandBuy struct {
	username string
	amount   string
	stock    string
}

func (c commandBuy) process(transaction int) string {
	logUserEvent(serverName, transaction, "BUY", c.username, c.stock, c.amount)
	return buy(c.username, c.stock, c.amount, transaction)
}

type commandCancelBuy struct {
	username string
}

func (c commandCancelBuy) process(transaction int) string {
	logUserEvent(serverName, transaction, "CANCEL_BUY", c.username, "", "")
	return cancelBuy(c.username, transaction)
}

type commandCommitBuy struct {
	username string
}

func (c commandCommitBuy) process(transaction int) string {
	logUserEvent(serverName, transaction, "COMMIT_BUY", c.username, "", "")
	return commitBuy(c.username, transaction)
}

func buy(userID string, stock string, pendingCashString string, transactionNum int) string {
	if !ratdatabase.UserExists(userID) {
		log.Println("user " + userID + " does not exist")
		return "No user exists"
	}

	pendingTransactionCash := stringToCents(pendingCashString)
	stockValue := getQuote(userID, stock, transactionNum)
	if stockValue <= 0 {
		return "Stock is worthless. I forbid you from buying"
	}
	stockAmount := pendingTransactionCash / stockValue
	currentBalance := ratdatabase.GetUserBalance(userID)

	if currentBalance < pendingTransactionCash {
		m := fmt.Sprintf("[%d] Not enough money for %s to perform buy", transactionNum, userID)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "BUY", userID, stock, pendingCashString, m)
		return m
	}

	if stockAmount == 0 {
		m := fmt.Sprintf("[%d] %s stock price(%d) higher than amount to purchase(%d)", transactionNum, stock, stockValue, pendingTransactionCash)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "BUY", userID, stock, pendingCashString, m)
		return m
	}

	//if has enough cash subtract and set aside from db
	newBalance := currentBalance - pendingTransactionCash
	ratdatabase.UpdateUserBalance(userID, newBalance)
	uuid := ratdatabase.InsertPendingBuyTransaction(userID, pendingTransactionCash, stock, stockValue)
	log.Printf("[%d] User %s buy transaction for %d %s@%.2f pending", transactionNum, userID, stockAmount, stock, float64(stockValue))

	go checkBuy(userID, uuid, transactionNum, pendingTransactionCash, stock)
	return "Buy for %s has been placed pending cancel/commit from user"
}

func checkBuy(userID string, uuid string, transactionNum int, pendingTransactionCash int, stock string) {
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

func cancelBuy(userID string, transactionNum int) string {
	uuid, holdingCash, stockName, _, exists := ratdatabase.GetLastPendingBuyTransaction(userID)

	if !exists {
		m := fmt.Sprintf("[%d] Cannot cancel buy. No buys pending", transactionNum)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "CANCEL_BUY", userID, "", "", m)
		return m
	}

	activeBalance := ratdatabase.GetUserBalance(userID)
	newBalance := activeBalance + holdingCash
	ratdatabase.UpdateUserBalance(userID, newBalance)

	//delete pending transaction
	ratdatabase.DeletePendingBuyTransaction(userID, uuid)
	m := fmt.Sprintf("[%d] Buy for stock:%s cancelled by user.", transactionNum, stockName)
	log.Printf(m)
	return m
}

func commitBuy(userID string, transactionNum int) string {
	uuid, holdingCash, stockName, stockPrice, exists := ratdatabase.GetLastPendingBuyTransaction(userID)

	if !exists {
		m := fmt.Sprintf("[%d] Cannot commit buy for %s. No buy pending", transactionNum, userID)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "COMMIT_BUY", userID, "", "", "Cannot commit buy. No buy pending")
		return m
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
	m := fmt.Sprintf("[%d] User %s now has %d more of stock %s", transactionNum, userID, stockBought, stockName)
	log.Println(m)

	if surplusCash != 0 {
		currentBalance := ratdatabase.GetUserBalance(userID)
		newBalance := currentBalance + surplusCash
		ratdatabase.UpdateUserBalance(userID, newBalance)
	}

	ratdatabase.DeletePendingBuyTransaction(userID, uuid)
	return m
}

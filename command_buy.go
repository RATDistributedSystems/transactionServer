package main

import (
	"fmt"
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

	fmt.Printf("Current: %d\nPending: %d\n", currentBalance, pendingTransactionCash)
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

	log.Printf("[%d] Cancelling '%s' request to buy %.2f of stock %s\n", transactionNum, userId, pendingTransactionCash, stock)
	ratdatabase.DeletePendingBuyTransaction(userId, uuid)

	// Returns users cash being held
	newerBalance := ratdatabase.GetUserBalance(userId)
	newererBalance := pendingTransactionCash + newerBalance
	ratdatabase.UpdateUserBalance(userId, newererBalance)
}

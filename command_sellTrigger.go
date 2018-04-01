package main

import (
	"fmt"
	"log"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandSetSellAmount struct {
	username string
	amount   string
	stock    string
}

func (c commandSetSellAmount) process(transaction int) string {
	logUserEvent(serverName, transaction, "SET_SELL_AMOUNT", c.username, c.stock, c.amount)
	return setSellAmount(c.username, c.stock, c.amount, transaction)
}

type commandSetSellTrigger struct {
	username string
	amount   string
	stock    string
}

func (c commandSetSellTrigger) process(transaction int) string {
	logUserEvent(serverName, transaction, "SET_SELL_TRIGGER", c.username, c.stock, c.amount)
	return setSellTrigger(c.username, c.stock, c.amount, transaction)
}

type commandCancelSetSell struct {
	username string
	stock    string
}

func (c commandCancelSetSell) process(transaction int) string {
	logUserEvent(serverName, transaction, "CANCEL_SET_SELL", c.username, c.stock, "")
	return cancelSellTrigger(c.username, c.stock, transaction)
}

func setSellAmount(userID string, stock string, pendingCashString string, transactionNum int) string {
	stockAmountToSell := stringToCents(pendingCashString)
	//check if user owns stock

	_, stockAmountOwned, ownsStock := ratdatabase.GetStockAmountOwned(userID, stock)

	if !ownsStock || stockAmountOwned == 0 {
		msg := "[%d] Not enough of stock %s (%d) for SellSetAmount %d for %s"
		m := fmt.Sprintf(msg, transactionNum, stock, stockAmountOwned, stockAmountToSell, userID)
		log.Printf(m)
		logErrorEvent(serverName, transactionNum, "SET_SELL_AMOUNT", userID, stock, pendingCashString, m)
		return m
	}

	currentStockPrice := getQuote(userID, stock, transactionNum)

	if currentStockPrice > stockAmountToSell {
		msg := "[%d] Current stock price (%d) is greater than amount wanting to be sold (%d)"
		m := fmt.Sprintf(msg, transactionNum, currentStockPrice, stockAmountToSell)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "SET_SELL_AMOUNT", userID, stock, pendingCashString, "Current stock price is greater than amount wanting to be sold.")
		return m
	}

	sellStockAmount := stockAmountToSell / currentStockPrice

	if stockAmountOwned < sellStockAmount {
		msg := "[%d] User %s does not own enough %s stock (%d) to sell %d amount"
		m := fmt.Sprintf(msg, transactionNum, userID, stock, stockAmountOwned, sellStockAmount)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "SET_SELL_AMOUNT", userID, stock, pendingCashString, m)
		return m
	}

	// Insert stock amount and return stock if it was already made
	oldSellStockAmount := ratdatabase.InsertSetSellTrigger(userID, stock, sellStockAmount)
	remainingStockAmount := stockAmountOwned + oldSellStockAmount - sellStockAmount
	ratdatabase.UpdateUserStockByUserAndStock(userID, stock, remainingStockAmount)

	return fmt.Sprintf("Successfully set SET_SELL_AMOUNT (%s) for %s", pendingCashString, stock)
}

func setSellTrigger(userID string, stock string, stockSellPrice string, transactionNum int) string {
	stockValueCents := stringToCents(stockSellPrice)
	stockAmountSet := ratdatabase.UpdateSellTriggerPrice(userID, stock, stockValueCents)

	if !stockAmountSet {
		msg := "[%d] User %s hasn't set stock amount for stock %s"
		m := fmt.Sprintf(msg, transactionNum, userID, stock)
		log.Println(m)
		logErrorEvent(serverName, transactionNum, "SET_SELL_TRIGGER", userID, stock, stockSellPrice, m)
		return m
	}

	return fmt.Sprintf("Successfully set SET_SELL_TRIGGER (%s) for %s", stockSellPrice, stock)
}

func cancelSellTrigger(userID string, stock string, transactionNum int) string {
	returnAmount := ratdatabase.CancelSellTrigger(userID, stock)

	if returnAmount == 0 {
		m := fmt.Sprintf("No trigger for %s set", stock)
		logErrorEvent(serverName, transactionNum, "CANCEL_SET_SELL", userID, stock, "", m)
		return m
	}

	_, oldStockAmount, _ := ratdatabase.GetStockAmountOwned(userID, stock)
	newStockAmount := oldStockAmount + returnAmount
	ratdatabase.UpdateUserStockByUserAndStock(userID, stock, newStockAmount)

	return fmt.Sprintf("Cancelled trigger for %s", stock)
}

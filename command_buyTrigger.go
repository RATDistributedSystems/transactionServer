package main

import (
	"fmt"
	"log"

	"github.com/RATDistributedSystems/utilities/ratdatabase"
)

type commandSetBuyAmount struct {
	username string
	amount   string
	stock    string
}

func (c commandSetBuyAmount) process(transaction int) string {
	logUserEvent(serverName, transaction, "SET_BUY_AMOUNT", c.username, c.stock, c.amount)
	return setBuyAmount(c.username, c.stock, c.amount, transaction)
}

type commandSetBuyTrigger struct {
	username string
	amount   string
	stock    string
}

func (c commandSetBuyTrigger) process(transaction int) string {
	logUserEvent(serverName, transaction, "SET_BUY_TRIGGER", c.username, c.stock, c.amount)
	return setBuyTrigger(c.username, c.stock, c.amount, transaction)
}

type commandCancelSetBuy struct {
	username string
	stock    string
}

func (c commandCancelSetBuy) process(transaction int) string {
	logUserEvent(serverName, transaction, "CANCEL_SET_BUY", c.username, c.stock, "")
	return cancelBuyTrigger(c.username, c.stock, transaction)
}

func setBuyAmount(userID string, stock string, pendingCashString string, transactionNum int) string {
	buyAmount := stringToCents(pendingCashString)
	userBalance := ratdatabase.GetUserBalance(userID)

	//if the user doesnt have enough funds cancel the allocation
	if userBalance < buyAmount {
		msg := "[%d] Not enough cash (%.2f) for buy amount trigger (%.2f) for %s"
		m := fmt.Sprintf(msg, transactionNum, float64(userBalance)/100, float64(buyAmount)/100, userID)
		log.Printf(m)
		logErrorEvent(serverName, transactionNum, "SET_BUY_AMOUNT", userID, stock, pendingCashString, m)
		return m
	}

	// Update user balance and add trigger
	oldTriggerAmount := ratdatabase.InsertSetBuyTrigger(userID, stock, buyAmount)
	newBalance := userBalance + oldTriggerAmount - buyAmount
	ratdatabase.UpdateUserBalance(userID, newBalance)

	return fmt.Sprintf("Sucessfully set SET_BUY_AMOUNT (%s) for %s", pendingCashString, stock)
}

func setBuyTrigger(userID string, stock string, stockPriceTriggerString string, transactionNum int) string {
	stockPriceTrigger := stringToCents(stockPriceTriggerString)
	buySetAmountExists := ratdatabase.UpdateBuyTriggerPrice(userID, stock, stockPriceTrigger)

	if !buySetAmountExists {
		msg := "[%d] Cannot set buy trigger price (%s). %s hasn't called BuySetAmount for stock %s"
		m := fmt.Sprintf(msg, transactionNum, stockPriceTriggerString, userID, stock)
		log.Printf(m)
		logErrorEvent(serverName, transactionNum, "SET_BUY_TRIGGER", userID, stock, stockPriceTriggerString, m)
		return m
	}

	return fmt.Sprintf("Successfully set SET_BUY_TRIGGER (%s) for %s", stockPriceTriggerString, stock)
}

func cancelBuyTrigger(userID string, stock string, transactionNum int) string {
	returnAmount := ratdatabase.CancelBuyTrigger(userID, stock)

	if returnAmount == 0 {
		m := fmt.Sprintf("No trigger for %s to cancel", stock)
		logErrorEvent(serverName, transactionNum, "CANCEL_SET_BUY", userID, stock, "", m)
		return m
	}

	userBalance := ratdatabase.GetUserBalance(userID)
	newBalance := userBalance + returnAmount
	ratdatabase.UpdateUserBalance(userID, newBalance)
	return fmt.Sprintf("Cancelled Buy Trigger for %s", stock)
}

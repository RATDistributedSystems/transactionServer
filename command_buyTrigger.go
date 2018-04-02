package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
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
	checkBuyTrigger(userID, stock, stockPriceTrigger, transactionNum)
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

func checkBuyTrigger(userId string, stock string, stockPriceTrigger int, transactionNum int) {
	operation := true
	for {
		//check the quote server every 5 seconds
		timer1 := time.NewTimer(time.Second * 1)
		<-timer1.C
		//if the trigger doesnt exist exit
		exists := checkTriggerExists(userId, stock, operation)
		if exists == false {
			return
		}
		currentStockPrice := getQuote(userId, stock, transactionNum)
		//execute the buy instantly if trigger condition is true
		if currentStockPrice <= stockPriceTrigger {
			var usableCash int
			var pendingCash int
			stockValue := currentStockPrice
			var remainingCash int
			var usid string
			var ownedstockname string
			var stockamount int
			var hasStock bool
			//check if user currently owns any of this stock
			iter := sessionGlobal.Query("SELECT usid, stock, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
			for iter.Scan(&usid, &ownedstockname, &stockamount) {
				if ownedstockname == stock {
					hasStock = true
					break
				}
			}
			if err := iter.Close(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}
			//If the user has some stock, add it to currently owned
			if hasStock == true {
				//grab pendingCash for the buy trigger
				if err := sessionGlobal.Query("SELECT pendingCash FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&pendingCash); err != nil {
					return
					//panic(fmt.Sprintf("problem getting usable cash form users", err))
				}
				//calculate amount of stocks can be bought
				buyableStocks := pendingCash / stockValue
				buyableStocks = buyableStocks + stockamount
				//remaining money
				remainingCash = pendingCash - (buyableStocks * stockValue)
				buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)
				//if the trigger doesnt exist exit
				exists := checkTriggerExists(userId, stock, operation)
				if exists == false {
					return
				}
				//insert new stock record
				if err := sessionGlobal.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE usid=" + usid + "").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				//check users available cash
				if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				//add available cash to leftover cash
				usableCash = usableCash + remainingCash
				usableCashString := strconv.FormatInt(int64(usableCash), 10)
				//re input the new cash value in to the user db
				if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				if err := sessionGlobal.Query("DELETE FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				return
			} else {
				//get pending cash in the trigger
				if err := sessionGlobal.Query("SELECT pendingCash FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&pendingCash); err != nil {
					return
					//panic(fmt.Sprintf("problem getting usable cash form users", err))
				}
				//IF USE DOESNT OWN ANY OF THIS STOCK
				//calculate amount of stocks can be bought
				if stockValue == 0 {
					return
				}
				buyableStocks := pendingCash / stockValue
				remainingCash = pendingCash - (buyableStocks * stockValue)
				buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)
				exists := checkTriggerExists(userId, stock, operation)
				if exists == false {
					return
				}
				//insert new stock record
				if err := sessionGlobal.Query("INSERT INTO userstocks (usid, userid, stockamount, stock) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + stock + "')").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				//check users available cash
				if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				//add available cash to leftover cash
				usableCash = usableCash + remainingCash
				usableCashString := strconv.FormatInt(int64(usableCash), 10)
				//re input the new cash value in to the user db
				if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				if err := sessionGlobal.Query("DELETE FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
					panic(fmt.Sprintf("problem creating session", err))
				}
				return
			}
		}
	}
}

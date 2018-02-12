package main

import (
	"fmt"
	"strconv"
	"time"
)

//Set maxmimum price of a stock before the stock gets auto bought
func setBuyTrigger(userId string, stock string, stockPriceTriggerString string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "SET_BUY_TRIGGER", userId, stock, stockPriceTriggerString)
	//convert trigger price from string to int cents
	stockPriceTrigger := stringToCents(stockPriceTriggerString)
	//fmt.Println(stockPriceTrigger);

	stockPriceTriggerString = strconv.FormatInt(int64(stockPriceTrigger), 10)

	//set the triggerValue and create thread to check the quote server
	if err := sessionGlobal.Query("UPDATE buyTriggers SET triggerValue =" + stockPriceTriggerString + " WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem setting trigger", err))
	}

	go checkBuyTrigger(userId, stock, stockPriceTrigger, transactionNum)

}

func checkBuyTrigger(userId string, stock string, stockPriceTrigger int, transactionNum int) {

	operation := true

	for {
		//check the quote server every 5 seconds
		timer1 := time.NewTimer(time.Second * 5)
		<-timer1.C

		//if the trigger doesnt exist exit

		exists := checkTriggerExists(userId, stock, operation)
		if exists == false {
			return
		}

		message := quoteRequest(userId, stock, transactionNum)
		currentStockPrice := stringToCents(message[0])

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
			iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
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
					panic(fmt.Sprintf("problem getting usable cash form users", err))
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
					panic(fmt.Sprintf("problem getting usable cash form users", err))
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
				if err := sessionGlobal.Query("INSERT INTO userstocks (usid, userid, stockamount, stockname) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + stock + "')").Exec(); err != nil {
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

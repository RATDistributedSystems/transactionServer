package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func commitBuy(userId string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "COMMIT_BUY", userId, "", "")
	buyExists := checkDependency("COMMIT_BUY", userId, "none")
	if buyExists == false {
		log.Printf("Cannot commit buy for %s #%d. No buy pending", userId, transactionNum)
		return
	}

	var pendingCash int
	var buyingstockName string
	var stockValue int
	var buyableStocks int
	var remainingCash int
	var usableCash int
	var uuid string
	//userId = strings.TrimSuffix(userId, "\n")

	if err := sessionGlobal.Query("select pid, stock, stockValue, pendingCash from buypendingtransactions where userId='"+userId+"'").Scan(&uuid, &buyingstockName, &stockValue, &pendingCash); err != nil {
		panic(err)
	}

	var usid string
	var ownedstockname string
	var stockamount int
	var hasStock bool

	//check if user currently owns any of this stock
	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
	for iter.Scan(&usid, &ownedstockname, &stockamount) {
		if ownedstockname == buyingstockName {
			hasStock = true
			break
		}
		//fmt.Println("STOCKS: ", stockname, stockamount)
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if hasStock == true {

		//calculate amount of stocks can be bought
		buyableStocks = pendingCash / stockValue
		buyableStocks = buyableStocks + stockamount
		//remaining money
		remainingCash = pendingCash - (buyableStocks * stockValue)

		buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)

		//update new stock record
		if err := sessionGlobal.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE usid=" + usid).Exec(); err != nil {
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

	} else {
		//IF USE DOESNT OWN ANY OF THIS STOCK
		//calculate amount of stocks can be bought
		buyableStocks = pendingCash / stockValue
		//remaining money
		remainingCash = pendingCash - (buyableStocks * stockValue)

		buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)

		userId = strings.TrimSuffix(userId, "\n")
		//insert new stock record
		if err := sessionGlobal.Query("INSERT INTO userstocks (usid, userid, stockamount, stockname) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + buyingstockName + "')").Exec(); err != nil {
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

	}

	//delete the pending transcation
	if err := sessionGlobal.Query("delete from buypendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}

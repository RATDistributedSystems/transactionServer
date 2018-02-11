package main

import (
	"fmt"
	"strconv"
	"strings"
)

func cancelSell(userId string, transactionNum int) {

	var uuid string
	var pendingCash int
	var pendingStock int
	var stock string
	var usid string
	var stockname string
	var stockamount int
	var totalStocks int
	var stocks int

	userId = strings.TrimSuffix(userId, "\n")

	sellExists := checkDependency("CANCEL_SELL", userId, "none")
	if sellExists == false {
		fmt.Println("cannot CANCEL SELL, no sell pending")
		return
	}

	if err := sessionGlobal.Query("select pid, userId, pendingcash, stock, stockvalue from sellpendingtransactions where userId='"+userId+"'"+" LIMIT 1").Scan(&uuid, &userId, &pendingCash, &stock, &pendingStock); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//get current user stocks
	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
	for iter.Scan(&usid, &stockname, &stockamount) {
		if stockname == stock {
			break

		}
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//convert stock value to stock amount
	stocks = pendingCash / pendingStock
	totalStocks = stocks + stockamount
	totalStocksString := strconv.FormatInt(int64(totalStocks), 10)

	//return user stocks
	if err := sessionGlobal.Query("UPDATE userstocks SET stockamount =" + totalStocksString + " WHERE usid=" + usid).Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete pending transaction
	if err := sessionGlobal.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}

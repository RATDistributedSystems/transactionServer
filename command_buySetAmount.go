package main

import (
	"fmt"
	"strconv"

	"github.com/twinj/uuid"
)

//sets aside the amount of money user wants to spend on a given stock
//executes prior to setTriggerValue
func setBuyAmount(userId string, stock string, pendingCashString string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "SET_BUY_AMOUNT", userId, stock, pendingCashString)
	var usableCash int
	pendingCash := stringToCents(pendingCashString)

	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem getting usable cash form users", err))
	}

	//if the user doesnt have enough funds cancel the allocation
	if usableCash < pendingCash {
		fmt.Println("Not enough money for this transaction")

		return
	}

	//allocate cash after being verified
	usableCash = usableCash - pendingCash
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString = strconv.FormatInt(int64(pendingCash), 10)
	//fmt.Println("Available Cash is greater than buy amount")
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem getting allocating user funds", err))
	}
	fmt.Println("Cash allocated")

	//generate UUID to input as a unique value
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	if err := sessionGlobal.Query("INSERT INTO buyTriggers (tid, pendingCash, stock, userid) VALUES (" + f + ", " + pendingCashString + ", '" + stock + "', '" + userId + "')").Exec(); err != nil {
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

}

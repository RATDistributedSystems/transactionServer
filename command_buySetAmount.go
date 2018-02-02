package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	//"strings"
	"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
	//"log"
)

//sets aside the amount of money user wants to spend on a given stock
//executes prior to setTriggerValue
func setBuyAmount(userId string, stock string, pendingCashString string,transactionNum int){

	//create session with cass database

	//Verify that use funds is greater than amount attempting to spend


	//Usable Cash is stored as cents
	var usableCash int


	//convert pendingCash from string to int of cents
	pendingCash := stringToCents(pendingCashString)

	/*
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//transactionNum_user += 1
	//transactionNum_user_string := strconv.FormatInt(int64(transactionNum_user), 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "SET_BUY_AMOUNT", userId, stock, pendingCashString)
	*/


	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem getting usable cash form users", err))
	}

	//Verify the pending cash vs the usable cash
	fmt.Println("\n" + userId)
	fmt.Println("usableCash")
	fmt.Println(usableCash)
	fmt.Println("pendingCash")
	fmt.Println(pendingCash)

	//if the user doesnt have enough funds cancel the allocation
	if usableCash < pendingCash{
		fmt.Println("Not enough money for this transaction")
		
		return
	}

	//allocate cash after being verified
	usableCash = usableCash - pendingCash;
	usableCashString := strconv.FormatInt(int64(usableCash), 10)
	pendingCashString = strconv.FormatInt(int64(pendingCash), 10)
	fmt.Println("Available Cash is greater than buy amount")
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem getting allocating user funds", err))
	}
	fmt.Println("Cash allocated")

	//Create an entry in the "Triggers" table to keep track of the initial buy amount setting

	//generate UUID to input as a unique value
	u := uuid.NewV4()
	f := uuid.Formatter(u, uuid.FormatCanonical)

	//buy operation flag
	//var operation string = "true"



	if err := sessionGlobal.Query("INSERT INTO buyTriggers (tid, pendingCash, stock, userid) VALUES (" + f + ", " + pendingCashString + ", '" + stock + "', '" + userId + "')").Exec(); err != nil{
		panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
	}

	

}
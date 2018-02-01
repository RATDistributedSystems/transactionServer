package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	"strings"
	//"github.com/twinj/uuid"
	//"time"
	//"github.com/go-redis/redis"
	//"log"
)

func processCommand(text string) []string{
	result := strings.Split(text, ",")
	for i := range result {
		fmt.Println(result[i])
	}
	return result;
}


func getUsableCash(userId string) int{
	var usableCash int
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}

	return usableCash
}



func stringToCents(x string) int{

	var dollars int
	var cents int

	fmt.Println(x)
	result := strings.Split(x, ".")
	for i := range result {
		fmt.Println(result[i])
	}

	if i, err := strconv.Atoi(result[0]); err == nil {
		dollars = i
		fmt.Println("dollar converted to int")
		fmt.Println(i)
	}

	result[1] = strings.TrimSuffix(result[1], "\n")
	if e, err := strconv.Atoi(result[1]); err == nil {
		cents = e
		fmt.Println("cents converted to int")
		fmt.Println(e)
	}

	dollars = dollars * 100
	var money int = dollars + cents

	return money
}

//chekcs if the command can be executed
//ie check if buy set before commit etc
func checkDependency(command string, userId string, stock string) bool{

	var count int

	if command == "COMMIT_BUY"{
		//check if a buy entry exists in buypendingtransactions table

		if err := sessionGlobal.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
		if count == 0{
			return false

		}else{

			return true
		}
	}
	if command == "COMMIT_SELL"{
		//check if a sell entry exists in sellpendingtransactions table
			//check if a sell entry exists in buypendingtransactions table

			if err := sessionGlobal.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}
			if count == 0{
				return false
			}else{
				return true
			}
	}
	if command == "CANCEL_BUY"{
		//check if a buy entry exists in buypendingtransactions table
		if err := sessionGlobal.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}

		
	}
	if command == "CANCEL_SELL"{
		//check if a sell entry exists in sellpendingtransactions table
		if err := sessionGlobal.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}
		
	}
	if command == "CANCEL_SET_BUY"{
		if err := sessionGlobal.Query("SELECT count(*) FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}
	}
	if command == "CANCEL_SET_SELL"{
		if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
		if count == 0{
			return false
		}else{
			return true
		}
		
	}

	return false
}

func addFunds(userId string, addCashAmount int){

	var usableCash int

	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	totalCash := usableCash + addCashAmount;
	totalCashString := strconv.FormatInt(int64(totalCash), 10)

	//return add funds to user
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

}








//check if the trigger hasn't been cancelled
func checkTriggerExists(userId string, stock string, operation bool) bool{


	var count int

	if operation == true {
		if err := sessionGlobal.Query("SELECT count(*) FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
	}else{
		if err := sessionGlobal.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil{
			panic(fmt.Sprintf("Problem inputting to buyTriggers Table", err))
		}
	}

	if count == 1 {
		return true
	}else{
		return false
	}
}









func checkStockOwnership(userId string, stock string) (int, string){

	var ownedstockname string
	var ownedstockamount int
	var usid string
	//var hasStock bool

	iter := sessionGlobal.Query("SELECT usid, stockname, stockamount FROM userstocks WHERE userid='"+ userId + "'").Iter()
	for iter.Scan(&usid, &ownedstockname, &ownedstockamount) {
		if (ownedstockname == stock){
			//hasStock = true
			break;
		}
	}
	if err := iter.Close(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//returns 0 if stock is empty
	return ownedstockamount, usid
	
}


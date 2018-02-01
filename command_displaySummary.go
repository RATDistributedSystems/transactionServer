package main

import (
	//"net"
	"fmt"
	//"bufio"
	//"os"
	//"github.com/gocql/gocql"
	"strconv"
	//"strings"
	//"github.com/twinj/uuid"
	"time"
	//"github.com/go-redis/redis"
	//"log"
)
func displaySummary(userId string, transactionNum int){
	//return user summary of their stocks, cash, triggers, etc
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	//check users available cash
	var usableCashString string
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCashString); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
	}
	transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
	logAccountTransactionEvent(timestamp_command, "TS1", transactionNum_string, "DISPLAY_SUMMARY", userId, usableCashString)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "DISPLAY_SUMMARY", userId, "", usableCashString)
}

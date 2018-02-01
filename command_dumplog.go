package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func dumpUser(userId string, filename string, transactionNum int) {
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum), 10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "DUMPLOG", "-1", "", "")
	log.Printf("Dumping log data for %s\n", userId)
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG, %s, %s", userId, filename))
}

func dump(filename string, transactionNum int) {
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	timestamp_command := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum), 10)
	logUserEvent(timestamp_command, "TS1", transactionNum_string, "DUMPLOG", "-1", "", "")
	log.Println("Dumping all log data")
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG, %s", filename))
}

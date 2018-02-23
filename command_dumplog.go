package main

import (
	"fmt"
	"log"
)

func dumpUser(userId string, filename string, transactionNum int) {
	//logUserEvent("TS1", transactionNum, "DUMPLOG", userId, "", "")
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s,%s", userId, filename))
	log.Printf("Dumping log data for %s\n", userId)
}

func dump(filename string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "DUMPLOG", "-1", "", "")
	log.Println("Dumping all log data")
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s", filename))
}

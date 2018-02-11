package main

import (
	"fmt"
	"log"
)

func dumpUser(userId string, filename string, transactionNum int) {
	log.Printf("Dumping log data for %s\n", userId)
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s,%s", userId, filename))
}

func dump(filename string, transactionNum int) {
	log.Println("Dumping all log data")
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s", filename))
}

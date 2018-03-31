package main

import (
	"fmt"
	"log"
)

type commandDumplog struct {
	username string
	filename string
}

func (c commandDumplog) process(transaction int) {
	logUserEvent(serverName, transaction, "DUMPLOG", c.username, "", "")
	if c.username == "-1" {
		dump(c.username, transaction)
	} else {
		dumpUser(c.username, c.filename, transaction)
	}

}

func dumpUser(userID string, filename string, transactionNum int) {
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s,%s", userID, filename))
	log.Printf("Dumping log data for %s\n", userID)
}

func dump(filename string, transactionNum int) {
	log.Println("Dumping all log data")
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s", filename))
}

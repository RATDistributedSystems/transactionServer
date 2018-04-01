package main

import (
	"fmt"
	"log"
)

type commandDumplog struct {
	username string
	filename string
}

func (c commandDumplog) process(transaction int) string {
	logUserEvent(serverName, transaction, "DUMPLOG", c.username, "", "")
	if c.username == "-1" {
		dump(c.username, transaction)
		return "" // only called from generator, not UI so who cares
	}

	return dumpUser(c.username, c.filename, transaction)
}

func dumpUser(userID string, filename string, transactionNum int) string {
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s,%s", userID, filename))
	m := fmt.Sprintf("Dumping log data for %s\n", userID)
	log.Printf(m)
	return m
}

func dump(filename string, transactionNum int) {
	log.Println("Dumping all log data")
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s", filename))
}

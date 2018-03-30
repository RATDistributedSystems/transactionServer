package main

import (
	"fmt"
	"log"
)

type commandDumplog struct {
	username string
}

func (c commandDumplog) process(transaction int) {
	logUserEvent(serverName, transaction, "DUMPLOG", "-1", "", "")
	dump(c.username, transaction)
}

type commandDumplogUsername struct {
	username string
	filename string
}

func (c commandDumplogUsername) process(transaction int) {
	logUserEvent(serverName, transaction, "DUMPLOG", c.username, "", "")
	dumpUser(c.username, c.filename, transaction)
}

func dumpUser(userId string, filename string, transactionNum int) {
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s,%s", userId, filename))
	log.Printf("Dumping log data for %s\n", userId)
}

func dump(filename string, transactionNum int) {
	log.Println("Dumping all log data")
	sendMsgToAuditServer(fmt.Sprintf("DUMPLOG,%s", filename))
}

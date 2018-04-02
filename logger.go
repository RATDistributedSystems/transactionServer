package main

import (
	"fmt"

	"github.com/RATDistributedSystems/utilities"
)

func sendMsgToAuditServer(msg string) {
	conn := auditPool.getConnection()
	fmt.Fprintln(conn, msg)
	auditPool.returnConnection(conn)
}

func logUserEvent(server string, transactionNum int, command string, userid string, stockSymbol string, funds string) {
	msg := fmt.Sprintf("User,%s,%s,%d,%s,%s,%s,%s", utilities.GetTimestamp(), server, transactionNum, command, userid, stockSymbol, funds)
	sendMsgToAuditServer(msg)
}

func logQuoteEvent(server string, transactionNum int, price string, stockSymbol string, userid string, quoteservertime string, cryptokey string) {
	msg := fmt.Sprintf("Quote,%s,%s,%d,%s,%s,%s,%s,%s", utilities.GetTimestamp(), server, transactionNum, price, stockSymbol, userid, quoteservertime, cryptokey)
	sendMsgToAuditServer(msg)
}

func logSystemEvent(server string, transactionNum int, command string, userid string, stockSymbol string, funds string) {
	msg := fmt.Sprintf("System,%s,%s,%d,%s,%s,%s,%s", utilities.GetTimestamp(), server, transactionNum, command, userid, stockSymbol, funds)
	sendMsgToAuditServer(msg)
}

func logAccountTransactionEvent(server string, transactionNum int, action string, userid string, funds string) {
	msg := fmt.Sprintf("Account,%s,%s,%d,%s,%s,%s", utilities.GetTimestamp(), server, transactionNum, action, userid, funds)
	sendMsgToAuditServer(msg)
}

func logErrorEvent(server string, transactionNum int, command string, userid string, stockSymbol string, funds string, err string) {
	msg := fmt.Sprintf("Error,%s,%s,%d,%s,%s,%s,%s,%s", utilities.GetTimestamp(), server, transactionNum, command, userid, stockSymbol, funds, err)
	sendMsgToAuditServer(msg)
}

func logDebugEvent(server string, transactionNum int, command string, userid string, stockSymbol string, funds string, debug string) {
	msg := fmt.Sprintf("Debug,%s,%s,%d,%s,%s,%s,%s,%s", utilities.GetTimestamp(), server, transactionNum, command, userid, stockSymbol, funds, debug)
	sendMsgToAuditServer(msg)
}

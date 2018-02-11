package main

import (
	"fmt"
	"strconv"
	"time"
)

func sendMsgToAuditServer(msg string) {
	conn := auditPool.getConnection()
	fmt.Println(conn, msg)
	auditPool.returnConnection(conn)
}

func logUserEvent(server string, transactionNum int, command string, userid string, stockSymbol string, funds string) {
	timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
	time := strconv.FormatInt(timestamp_time, 10)
	transactionNum_string := strconv.FormatInt(int64(transactionNum), 10)
	msg := fmt.Sprintf("User,%s,%s,%s,%s,%s,%s,%s", time, server, transactionNum_string, command, userid, stockSymbol, funds)
	sendMsgToAuditServer(msg)
}

func logQuoteEvent(time string, server string, transactionNum string, price string, stockSymbol string, userid string, quoteservertime string, cryptokey string) {
	msg := fmt.Sprintf("Quote,%s,%s,%s,%s,%s,%s,%s,%s", time, server, transactionNum, price, stockSymbol, userid, quoteservertime, cryptokey)
	sendMsgToAuditServer(msg)
}

func logSystemEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string) {
	msg := fmt.Sprintf("System,%s,%s,%s,%s,%s,%s,%s", time, server, transactionNum, command, userid, stockSymbol, funds)
	sendMsgToAuditServer(msg)
}

func logAccountTransactionEvent(time string, server string, transactionNum string, action string, userid string, funds string) {
	msg := fmt.Sprintf("Account,%s,%s,%s,%s,%s,%s", time, server, transactionNum, action, userid, funds)
	sendMsgToAuditServer(msg)
}

func logErrorEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, err string) {
	msg := fmt.Sprintf("User,%s,%s,%s,%s,%s,%s,%s,%s", time, server, transactionNum, command, userid, stockSymbol, funds, err)
	sendMsgToAuditServer(msg)
}

func logDebugEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, debug string) {
	msg := fmt.Sprintf("User,%s,%s,%s,%s,%s,%s,%s,%s", time, server, transactionNum, command, userid, stockSymbol, funds, debug)
	sendMsgToAuditServer(msg)
}

package main


//import "net"
import "fmt"
//import "bufio"
//import "os"
//import "github.com/gocql/gocql"
//import "strconv"
//import "strings"
//import "github.com/go-redis/redis"
//import "github.com/twinj/uuid"
//import "time"
//import "log"


func sendMsgToAuditServer(msg string) {
	//addr, protocol := configurationServer.GetServerDetails("audit")
	//conn, err := net.Dial(protocol, addr)
	//if err != nil {
		//log.Fatalf("Couldn't Connect to Audit server: " + err.Error())
	//}
	//defer conn.Close()
	fmt.Fprintln(logConnection, msg)
}


func logUserEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string) {
	msg := fmt.Sprintf("User, %s, %s, %s, %s, %s, %s, %s", time, server, transactionNum, command, userid, stockSymbol, funds)
	sendMsgToAuditServer(msg)
}

func logQuoteEvent(time string, server string, transactionNum string, price string, stockSymbol string, userid string, quoteservertime string, cryptokey string) {
	msg := fmt.Sprintf("Quote, %s, %s, %s, %s, %s, %s, %s, %s", time, server, transactionNum, price, stockSymbol, userid, quoteservertime, cryptokey)
	sendMsgToAuditServer(msg)
}

func logSystemEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string) {
	msg := fmt.Sprintf("System, %s, %s, %s, %s, %s, %s, %s", time, server, transactionNum, command, userid, stockSymbol, funds)
	sendMsgToAuditServer(msg)
}

func logAccountTransactionEvent(time string, server string, transactionNum string, action string, userid string, funds string) {
	msg := fmt.Sprintf("Account, %s, %s, %s, %s, %s, %s", time, server, transactionNum, action, userid, funds)
	sendMsgToAuditServer(msg)
}

func logErrorEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, err string) {
	msg := fmt.Sprintf("User, %s, %s, %s, %s, %s, %s, %s, %s", time, server, transactionNum, command, userid, stockSymbol, funds, err)
	sendMsgToAuditServer(msg)
}

func logDebugEvent(time string, server string, transactionNum string, command string, userid string, stockSymbol string, funds string, debug string) {
	msg := fmt.Sprintf("User, %s, %s, %s, %s, %s, %s, %s, %s", time, server, transactionNum, command, userid, stockSymbol, funds, debug)
	sendMsgToAuditServer(msg)
}

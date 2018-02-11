package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func commitSell(userId string, transactionNum int) {

	var uuid string
	var pendingCash int
	var usableCash int
	var stock string
	userId = strings.TrimSuffix(userId, "\n")

	sellExists := checkDependency("COMMIT_SELL", userId, "none")
	if sellExists == false {
		log.Printf("Cannot commit sell for user:  %s #%d. No sells pending.", userId, transactionNum)
		return
	}

	if err := sessionGlobal.Query("select pid,stock from sellpendingtransactions where userid='"+userId+"'").Scan(&uuid, &stock); err != nil {
		panic(fmt.Sprintf("problem", err))
	}

	//get pending cash to be added to user account
	if err := sessionGlobal.Query("select pid, pendingcash from sellpendingtransactions where userid='"+userId+"'").Scan(&uuid, &pendingCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	if uuid == "" {
		return
	}

	//get current users cash
	if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//add available cash to leftover cash
	usableCash = usableCash + pendingCash
	usableCashString := strconv.FormatInt(int64(usableCash), 10)

	//re input the new cash value in to the user db
	if err := sessionGlobal.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}

	//delete the pending transcation
	if err := sessionGlobal.Query("delete from sellpendingtransactions where pid=" + uuid + " and userid='" + userId + "'").Exec(); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
}

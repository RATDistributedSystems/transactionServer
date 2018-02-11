package main

func cancelSellTrigger(userId string, stock string, transactionNum int) {
	logUserEvent("TS1", transactionNum, "CANCEL_SET_SELL", userId, stock, "")
	sellExists := checkDependency("CANCEL_SET_SELL", userId, stock)
	if sellExists == false {
		//fmt.Println("cannot CANCEL, no sells pending")
		return
	}

	if err := sessionGlobal.Query("DELETE FROM sellTriggers WHERE userId='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
		panic(err)
	}
}

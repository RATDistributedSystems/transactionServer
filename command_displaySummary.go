package main

func displaySummary(userId string, transactionNum int) {
	// return user summary of their stocks, cash, triggers, etc
	// Not implemented yet
	/*
		timestamp_time := (time.Now().UTC().UnixNano()) / 1000000
		timestamp_command := strconv.FormatInt(timestamp_time, 10)
		transactionNum_string := strconv.FormatInt(int64(transactionNum),10)
		//logAccountTransactionEvent(timestamp_command, "TS1", transactionNum_string, "DISPLAY_SUMMARY", userId, usableCashString)
		logUserEvent(timestamp_command, "TS1", transactionNum_string, "DISPLAY_SUMMARY", userId, "", "")

		//check users available cash
		//var usableCashString string
		if err := sessionGlobal.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCashString); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
		}
	*/

}

package main

type commandDisplaySummary struct {
	username string
}

func (c commandDisplaySummary) process(transaction int) {
	logUserEvent(serverName, transaction, "DISPLAY_SUMMARY", c.username, "", "")
	displaySummary(c.username, transaction)
}

func displaySummary(userId string, transactionNum int) {

	// return user summary of their stocks, cash, triggers, etc
	// Not implemented yet
}

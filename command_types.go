package main

type command interface {
	process(int)
}

type commandSetBuyAmount struct {
	username string
	amount   string
	stock    string
}

func (c commandSetBuyAmount) process(transaction int) {
	logUserEvent(serverName, transaction, "SET_BUY_AMOUNT", c.username, c.stock, c.amount)
	setBuyAmount(c.username, c.stock, c.amount, transaction)
}

type commandSetBuyTrigger struct {
	username string
	amount   string
	stock    string
}

func (c commandSetBuyTrigger) process(transaction int) {
	logUserEvent(serverName, transaction, "SET_BUY_TRIGGER", c.username, c.stock, c.amount)
	setBuyTrigger(c.username, c.stock, c.amount, transaction)
}

type commandCancelSetBuy struct {
	username string
	stock    string
}

func (c commandCancelSetBuy) process(transaction int) {
	logUserEvent(serverName, transaction, "CANCEL_SET_BUY", c.username, c.stock, "")
	cancelBuyTrigger(c.username, c.stock, transaction)
}

type commandSetSellAmount struct {
	username string
	amount   string
	stock    string
}

func (c commandSetSellAmount) process(transaction int) {
	logUserEvent(serverName, transaction, "SET_SELL_AMOUNT", c.username, c.stock, c.amount)
	setSellAmount(c.username, c.stock, c.amount, transaction)
}

type commandSetSellTrigger struct {
	username string
	amount   string
	stock    string
}

func (c commandSetSellTrigger) process(transaction int) {
	logUserEvent(serverName, transaction, "SET_SELL_TRIGGER", c.username, c.stock, c.amount)
	setSellTrigger(c.username, c.stock, c.amount, transaction)
}

type commandCancelSetSell struct {
	username string
	stock    string
}

func (c commandCancelSetSell) process(transaction int) {
	logUserEvent(serverName, transaction, "CANCEL_SET_SELL", c.username, c.stock, "")
	cancelSellTrigger(c.username, c.stock, transaction)
}

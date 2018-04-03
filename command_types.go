package main

import (
	"fmt"
	"strconv"
	"strings"
)

type command interface {
	process(int) string
}

type stockDetails struct {
	name   string
	amount int
}

type stockTriggerDetails struct {
	name   string
	amount int
	price  int
}

type userAcount struct {
	name         string
	balance      int
	stocks       []stockDetails
	buytriggers  []stockTriggerDetails
	selltriggers []stockTriggerDetails
}

func createTableHead(element ...string) string {
	middle := strings.Join(element, "</th><th>")
	return fmt.Sprintf("<thead><tr><th>%s</th></tr></thead>", middle)
}

func createTableBody(element ...string) string {
	middle := strings.Join(element, "</td><td>")
	return fmt.Sprintf("<tbody><tr><td>%s</td></tr></tbody>", middle)
}

func createTable(title string, head string, rows ...string) string {
	tableTemp := "<div class=\"container-fluid\"><h2>%s</h2><table class=\"table table-striped\">%s%s</table></div><br>"

	body := strings.Join(rows, " ")
	return fmt.Sprintf(tableTemp, title, head, body)
}

func (u userAcount) String() string {
	// Create table for user account details
	userHead := createTableHead("User", "Balance")
	userBody := createTableBody(u.name, centsToString(u.balance))
	userTable := createTable("User Account Summary", userHead, userBody)

	// table for stock owned details
	stockHead := createTableHead("Stock", "Amount")
	var stockBodies []string
	for _, stock := range u.stocks {
		row := createTableBody(stock.name, strconv.Itoa(stock.amount))
		stockBodies = append(stockBodies, row)
	}
	stockTable := createTable("Stocks Owned", stockHead, stockBodies...)

	// table for triggers
	buytriggerHead := createTableHead("Stock", "Trigger Amount", "Trigger Price")
	var buytriggerBodies []string
	for _, trigger := range u.buytriggers {
		row := createTableBody(trigger.name, strconv.Itoa(trigger.amount), centsToString(trigger.price))
		buytriggerBodies = append(buytriggerBodies, row)
	}
	buytriggerTable := createTable("Buy Triggers", buytriggerHead, buytriggerBodies...)

	selltriggerHead := createTableHead("Stock", "Trigger Amount", "Trigger Price")
	var selltriggerBodies []string
	for _, trigger := range u.buytriggers {
		row := createTableBody(trigger.name, strconv.Itoa(trigger.amount), centsToString(trigger.price))
		selltriggerBodies = append(selltriggerBodies, row)
	}
	sellriggerTable := createTable("Sell Triggers", selltriggerHead, selltriggerBodies...)

	return userTable + stockTable + buytriggerTable + sellriggerTable
}

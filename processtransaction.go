package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
)

type commandValidator struct {
	comm           command
	userRequired   bool
	stockRequired  bool
	amountRequired bool
}

func getPostInformation(r *http.Request) (c command, transaction int, e error) {
	r.ParseForm()
	f := r.PostForm
	commandType, _ := getParameter(f, "command")

	switch commandType {
	case "ADD":
		u, eu := getParameter(f, "username")
		a, ea := getParameter(f, "amount")
		c = commandAdd{u, a}
		if eu != nil || ea != nil {
			e = eu
		}
	case "BUY":
		u, eu := getParameter(f, "username")
		a, ea := getParameter(f, "amount")
		s, es := getParameter(f, "stock")
		c = commandBuy{u, a, s}
		if eu != nil || ea != nil || es != nil {
			e = eu
		}
	case "SELL":
		u, eu := getParameter(f, "username")
		a, ea := getParameter(f, "amount")
		s, es := getParameter(f, "stock")
		c = commandSell{u, a, s}
		if eu != nil || ea != nil || es != nil {
			e = eu
		}
	case "QUOTE":
		u, eu := getParameter(f, "username")
		s, es := getParameter(f, "stock")
		c = commandQuote{u, s}
		if eu != nil || es != nil {
			e = eu
		}
	case "COMMIT_BUY":
		u, eu := getParameter(f, "username")
		c = commandCommitBuy{u}
		if eu != nil {
			e = eu
		}
	case "COMMIT_SELL":
		u, eu := getParameter(f, "username")
		c = commandCommitSell{u}
		if eu != nil {
			e = eu
		}
	case "CANCEL_BUY":
		u, eu := getParameter(f, "username")
		c = commandCancelBuy{u}
		if eu != nil {
			e = eu
		}
	case "CANCEL_SELL":
		u, eu := getParameter(f, "username")
		c = commandCanceSell{u}
		if eu != nil {
			e = eu
		}
	case "SET_BUY_AMOUNT":
		u, eu := getParameter(f, "username")
		a, ea := getParameter(f, "amount")
		s, es := getParameter(f, "stock")
		c = commandSetBuyAmount{u, a, s}
		if eu != nil || ea != nil || es != nil {
			e = eu
		}
	case "SET_BUY_TRIGGER":
		u, eu := getParameter(f, "username")
		a, ea := getParameter(f, "amount")
		s, es := getParameter(f, "stock")
		c = commandSetBuyTrigger{u, a, s}
		if eu != nil || ea != nil || es != nil {
			e = eu
		}
	case "CANCEL_SET_BUY":
		u, eu := getParameter(f, "username")
		s, es := getParameter(f, "stock")
		c = commandCancelSetBuy{u, s}
		if eu != nil || es != nil {
			e = eu
		}
	case "SET_SELL_AMOUNT":
		u, eu := getParameter(f, "username")
		a, ea := getParameter(f, "amount")
		s, es := getParameter(f, "stock")
		c = commandSetSellAmount{u, a, s}
		if eu != nil || ea != nil || es != nil {
			e = eu
		}
	case "SET_SELL_TRIGGER":
		u, eu := getParameter(f, "username")
		a, ea := getParameter(f, "amount")
		s, es := getParameter(f, "stock")
		c = commandSetSellTrigger{u, a, s}
		if eu != nil || ea != nil || es != nil {
			e = eu
		}
	case "CANCEL_SET_SELL":
		u, eu := getParameter(f, "username")
		s, es := getParameter(f, "stock")
		c = commandCancelSetSell{u, s}
		if eu != nil || es != nil {
			e = eu
		}
	case "DUMPLOG":
		u, eu := getParameter(f, "username")
		s, es := getParameter(f, "stock")
		c = commandDumplog{u, s}
		if eu != nil || es != nil {
			e = eu
		}
	case "DISPLAY_SUMMARY":
		u, eu := getParameter(f, "username")
		c = commandDisplaySummary{u}
		if eu != nil {
			e = eu
		}
	default:
		e = errors.New("Invalid Command")
	}

	transaction = getTransactionNumber(f)
	return
}

func getTransactionNumber(f url.Values) (transaction int) {
	transactionString, et := getParameter(f, "transaction")
	if et != nil || transactionString == "0" {
		atomic.AddInt64(&__transaction_number, 1)
		transaction = int(atomic.LoadInt64(&__transaction_number))
		return
	}

	transaction, _ = strconv.Atoi(transactionString)
	return

}

func getParameter(f url.Values, p string) (string, error) {
	splice := f[p]
	if splice == nil || len(splice) != 1 || splice[0] == "" {
		return "", fmt.Errorf("Missing/Invalid Required Username")
	}
	return strings.TrimSpace(splice[0]), nil
}

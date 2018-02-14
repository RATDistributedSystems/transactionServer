package main

import (
	//"bufio"
	//"fmt"
	"strings"
)

func cacheExists(stock string) bool{
	return true

}

func cacheReturn(stock string) []string {
	//query the db to get the data
	message := "quote, here"
	//break comma delimitted data in to a message
	messageArray := strings.Split(message, ",")
	//return the quote array
	return messageArray
}
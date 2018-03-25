package main

import (
	//"bufio"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var redisConn *redis.Client

const ttl = time.Second * 60

func initRedis() {
	addr, proto := configurationServer.GetServerDetails("redis")
	redisConn = redis.NewClient(&redis.Options{
		Network:  proto,
		Addr:     addr,
		Password: "",
		DB:       0, //using the default DB
	})
	_, err := redisConn.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.Println("Redis Connection Created")
}

func cacheAdd(stockName string, stockInfo string) {
	fmt.Printf("Adding %s to cache \n", stockName)
	err := redisConn.Set(stockName, stockInfo, ttl).Err()
	if err != nil {
		panic(err)
	}
}

func cacheExists(stock string) (bool, int) {
	fmt.Printf("checking if %s exists", stock)
	val, err := redisConn.Get(stock).Result()

	//fmt.Println(val)
	if err == redis.Nil {
		fmt.Printf("Cache %s DNE \n", stock)
		return false, 0
	} else if err != nil {
		panic(err)
	} else {
		messageArray := strings.Split(val, ",")
		fmt.Printf("Cache %s Exists \n", stock)
		return true, stringToCents(messageArray[0])
	}

}

func cacheReturn(stock string) int {
	fmt.Printf("Returning Cache Key %s \n", stock)
	val, err := redisConn.Get(stock).Result()
	if err != nil {
		panic(err)
	}
	//break comma delimitted data in to a message
	messageArray := strings.Split(val, ",")
	//return the quote array
	//return messageArray
	return stringToCents(messageArray[0])
}

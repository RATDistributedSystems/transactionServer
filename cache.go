package main

import (
	//"bufio"

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
	err := redisConn.Set(stockName, stockInfo, ttl).Err()
	if err != nil {
		panic(err)
	}
}

func cacheExists(stock string) (bool, int) {
	val, err := redisConn.Get(stock).Result()

	//fmt.Println(val)
	if err == redis.Nil {
		return false, 0
	} else if err != nil {
		panic(err)
	} else {
		messageArray := strings.Split(val, ",")
		return true, stringToCents(messageArray[0])
	}

}

func cacheReturn(stock string) int {
	val, err := redisConn.Get(stock).Result()
	if err != nil {
		panic(err)
	}

	messageArray := strings.Split(val, ",")
	return stringToCents(messageArray[0])
}

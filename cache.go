package main

import (
	//"bufio"
	"fmt"
	"time"
	"strings"
	"github.com/go-redis/redis"
)

var redisConn *redis.Client

const ttl = time.Second * 60


func initRedis(){
	redisConn = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0, //using the default DB
	})
	fmt.Println("Redis Connection Created")
}

func cacheAdd(stockName string, stockInfo string){
	fmt.Printf("Adding %s to cache \n", stockName)
	err := redisConn.Set(stockName, stockInfo, ttl).Err()
	if err != nil {
		panic(err)
	}
}

func cacheExists(stock string) bool{
	fmt.Printf("checking if %s exists", stock)
	val, err := redisConn.Get(stock).Result()
	fmt.Println(val)
	if err == redis.Nil{
		fmt.Printf("Cache %s DNE \n", stock)
		return false
	}else if err != nil {
		panic(err)
	}else {
		fmt.Printf("Cache %s Exists \n", stock)
		return true
	}

}

func cacheReturn(stock string) int {
	fmt.Printf("Returning Cache Key %s \n", stock)
	val, err := redisConn.Get(stock).Result()
	if err != nil{
		panic(err)
	}
	//break comma delimitted data in to a message
	messageArray := strings.Split(val, ",")
	//return the quote array
	//return messageArray
	return stringToCents(messageArray[0])
}
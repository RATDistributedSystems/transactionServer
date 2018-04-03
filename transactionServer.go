package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/RATDistributedSystems/utilities"
	"github.com/RATDistributedSystems/utilities/ratdatabase"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"github.com/twinj/uuid"
)

var configurationServer = utilities.Load()
var serverName = gocql.TimeUUID().String()
var sessionGlobal *gocql.Session
var auditPool *connectionPool

var __transaction_number int64

func main() {


	initBenchFiles()
	uuid.Init()
	configurationServer.Pause()
	auditPool = initializePool(150, 190, "audit")
	initCassandra()
	initRedis()
	initHTTPServer()
}

func initHTTPServer() {
	addr, _ := configurationServer.GetListnerDetails("transaction")
	// Enable HTTP handlers
	router := httprouter.New()
	router.GET("/", getURL)
	router.GET("/add", getURL)
	router.GET("/buy", getURL)
	router.GET("/buytrigger", getURL)
	router.GET("/commit", getURL)
	router.GET("/quote", getURL)
	router.GET("/sell", getURL)
	router.GET("/selltrigger", getURL)
	router.GET("/summary", getURL)
	router.POST("/result", requestHandler)
	log.Printf("Serving on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func initCassandra() {
	//connect to database
	hostname := configurationServer.GetValue("transdb_ip")
	keyspace := configurationServer.GetValue("transdb_keyspace")
	protocol := configurationServer.GetValue("transdb_proto")
	ratdatabase.InitCassandraConnection(hostname, keyspace, protocol)
	sessionGlobal = ratdatabase.CassandraConnection
}

func requestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	comm, transaction, err := getPostInformation(r)
	if err != nil {
		errorResponse(w, err.Error())
		r.Body.Close()
		log.Println(err.Error())
		return
	}

	resp := comm.process(transaction)
	response(w, resp)
	r.Body.Close()

}

func initBenchFiles(){
	fmt.Println("Creating Benchmark Files")
	createFile("add.txt")
	createFile("addCQL.txt")
	createFile("buy.txt")
	createFile("buyCQL.txt")
	createFile("sell.txt")
	createFile("cancelBuy.txt")
	createFile("cancelBuyCQL.txt")
	createFile("cancelSell.txt")
	createFile("commitBuy.txt")
	createFile("commitSell.txt")
	createFile("quote.txt")
	createFile("quoteCache.txt")
	createFile("displaySummary.txt")
	createFile("dumpLog.txt")
	createFile("setBuyAmount.txt")
	createFile("setBuyTrigger.txt")
	createFile("cancelBuyTrigger.txt")
	createFile("setSellAmount.txt")
	createFile("setSellTrigger.txt")
	createFile("cancelSellTrigger.txt")

}



package main

import (
	"log"
	"net/http"
	"sync/atomic"

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

//var quotePool = initializePool(100, "quote")

func main() {

	uuid.Init()
	configurationServer.Pause()
	auditPool = initializePool(150, 190, "audit")
	initCassandra()
	initRedis()
	initHTTPServer()

}

func initHTTPServer() {
	addrWS, _ := configurationServer.GetListnerDetails("transaction")
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
	log.Printf("Serving on %s", addrWS)
	log.Fatal(http.ListenAndServe(addrWS, router))
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

	r.ParseForm()
	comm, err := getPostInformation(r.PostForm)
	if err != nil {
		ErrorResponse(w, err.Error())
		r.Body.Close()
		log.Println(err.Error())
		return
	}

	SuccessResponse(w)
	r.Body.Close()
	atomic.AddInt64(&__transaction_number, 1)
	transaction := atomic.LoadInt64(&__transaction_number)
	go comm.process(int(transaction))
}

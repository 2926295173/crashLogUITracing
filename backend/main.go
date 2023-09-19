package main

import (
	"clash-tracing/tracingstruct"
	"clash-tracing/utils"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed static/*
var f embed.FS

var (
	wg sync.WaitGroup
	db *sql.DB
)

func HandleTracing() {
	defer wg.Done()

	serverURL := "ws://127.0.0.1:9090/profile/tracing?token="

	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server: /profile/tracing")

	handlerMap := map[string]func([]byte, *sql.DB){
		"DNSRequest": utils.HandleDNSRequestLog,
		"RuleMatch":  utils.HandleRuleMatchLog,
		"ProxyDial":  utils.HandleProxyDialLog,
	}

	for {
		_, receivedMessage, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		logType := tracingstruct.LogType{}
		err = json.Unmarshal(receivedMessage, &logType)
		if err != nil {
			continue
		}

		if handlerMap[logType.Type] != nil {
			handlerMap[logType.Type](receivedMessage, db)
		}
	}
}

func HandleTraffic() {
	defer wg.Done()

	serverURL := "ws://127.0.0.1:9090/traffic?token="

	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server: /traffic")

	for {
		_, receivedMessage, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		utils.HandleTraffic(receivedMessage, db)
	}
}

func HandleHttpApi() {
	defer wg.Done()

	routers := utils.GetRouters()

	for _, router := range routers {
		for k, v := range router {
			http.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
				v(w, r, db, f)
			})
			fmt.Printf("Route [%s] Registered\n", k)
		}
	}

	fmt.Println("Api running @ http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

func main() {
	db = utils.ConnectDB()
	defer db.Close()

	utils.InitDBTables(db)

	wg.Add(1)
	go HandleTracing()

	wg.Add(1)
	go HandleTraffic()

	wg.Add(1)
	go HandleHttpApi()

	wg.Wait()
}

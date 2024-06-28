package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var goCount int64

func ws(w http.ResponseWriter, r *http.Request) {
	// upgrade http to ws
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade conn fail: ", err)
		return
	}

	n := atomic.AddInt64(&goCount, 1)
	if goCount%100 == 0 {
		log.Println("total number of conn = ", n)
	}
	defer func() {
		n := atomic.AddInt64(&goCount, -1)
		if goCount%100 == 0 {
			log.Println("total number of conn = ", n)
		}
		_ = conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read conn msg err ", err)
			return
		}
		log.Println(string(msg))
	}
}

func main() {

	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatal("pprof fail: ", err)
		}
	}()

	http.HandleFunc("/", ws)
	if err := http.ListenAndServe("localhost:8888", nil); err != nil {
		log.Fatal(err)
	}
}

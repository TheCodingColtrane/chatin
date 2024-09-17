package main

import (
	"flag"
	"log"
	"messenger/server"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	hub := server.NewHub()
	go hub.Run()
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

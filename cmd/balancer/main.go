package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"guthub.com/povqdlo/tz_balancer/server"
)

// Pool ...
var Pool server.ServerPool

//counterIn - Счетчик входящих сообщений
var CounterIn uint64

func init() {

	for _, c := range *server.ReadConfig() {

		newServer := server.NewServer(&c)
		server.Run(newServer)
		Pool = append(Pool, newServer)
	}
}

func main() {

	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			CounterIn++
			Pool.LoadBalance(&Pool, w, r)
		}
	})

	lbPort := "8080"
	go http.ListenAndServe(fmt.Sprintf(":%s", lbPort), nil)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	fmt.Printf("Количество входящих: %d\tКоличество обработанных: %d", CounterIn, server.CounterServer)

}

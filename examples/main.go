package main

import (
	"flag"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	listenAddr = flag.String("addr", ":8080", "listen address")
)

func main() {
	a := newApp()
	r := mux.NewRouter()
	r.HandleFunc("/", a.indexPage)
	r.HandleFunc("/users/", a.usersListPage)
	r.HandleFunc("/users/edit/{id}", a.usersEditPage)

	log.Info("Listening on: ", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, r); err != nil {
		log.Fatal(err)
	}
}

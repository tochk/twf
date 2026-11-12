package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// maxFormSize limits the size of posted forms, including uploaded files.
const maxFormSize = 32 << 20

var (
	listenAddr = flag.String("addr", ":8081", "listen address")
)

func main() {
	flag.Parse()

	a := newApp()
	r := mux.NewRouter()
	r.HandleFunc("/", a.indexPage)
	r.HandleFunc("/users/", a.usersListPage)
	r.HandleFunc("/users/edit/{id}", a.usersEditPage)

	slog.Info("listening", "addr", *listenAddr)
	srv := &http.Server{
		Addr:              *listenAddr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

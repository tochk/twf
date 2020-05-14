package main

import (
	"fmt"
	"net/http"
)

func (a *app) indexPage(w http.ResponseWriter, r *http.Request) {
	//todo
	fmt.Fprintln(w, "not implemented")
}

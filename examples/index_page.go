package main

import (
	"fmt"
	"net/http"
)

type loginData struct {
	Login    string `twf:"name:login,title:Login"`
	Password string `twf:"name:password,title:Password,type:password"`
}

func (a *app) indexPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var logData loginData
		data, err := a.twfUser.Form("Login", &logData, "")
		if err != nil {
			fmt.Fprint(w, "Err: ", err)
		}
		fmt.Fprint(w, data)
		return
	case http.MethodPost:

	}
}

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tochk/twf"
	"net/http"
)

// loginData structure for log in data - login and password
type loginData struct {
	// Login field in form - simple text field with name, title and placeholder described below
	Login string `twf:"name:login,title:Login,placeholder:Enter any login"`
	// Password field in form - password field with name, title and placeholder described below
	Password string `twf:"name:password,title:Password,type:password,placeholder:Enter any password"`
}

// indexPage handler for / page
func (a *app) indexPage(w http.ResponseWriter, r *http.Request) {
	var logData loginData

	switch r.Method {
	case http.MethodGet:
		// build add form by loginData structure
		data, err := a.twfUser.AddForm("Login", &logData, "")
		if err != nil {
			fmt.Fprint(w, "Err: ", err)
			return
		}

		// print form page to user
		fmt.Fprint(w, data)
	case http.MethodPost:
		// don't forget to parse form before calling twf.PostFormToStruct function
		r.ParseMultipartForm(32 << 20)

		// parse form
		err := twf.PostFormToStruct(&logData, r)
		if err != nil {
			fmt.Fprint(w, "Err: ", err)
			return
		}

		// print results
		log.Printf("entered login data: %+v", logData)

		// redirect to users page
		http.Redirect(w, r, "/users/", http.StatusFound)
	}
}

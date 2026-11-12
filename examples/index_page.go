package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tochk/twf"
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
	r.Body = http.MaxBytesReader(w, r.Body, maxFormSize)

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
		if err := r.ParseMultipartForm(maxFormSize); err != nil && !errors.Is(err, http.ErrNotMultipart) { //nolint:gosec // G120 false positive: body capped by http.MaxBytesReader
			fmt.Fprint(w, "Err: ", err)
			return
		}

		// parse form
		err := twf.PostFormToStruct(&logData, r)
		if err != nil {
			fmt.Fprint(w, "Err: ", err)
			return
		}

		// print results
		slog.Info("entered login data", "data", logData)

		// redirect to users page
		http.Redirect(w, r, "/users/", http.StatusFound)
	}
}

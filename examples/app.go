package main

import (
	"github.com/tochk/twf"
	"github.com/tochk/twf/examples/example_templates"
	"github.com/tochk/twf/twftemplates"
)

// app base app structure with user and admin twf instances
type app struct {
	twfUser  *twf.TWF
	twfAdmin *twf.TWF
}

// newApp function for creating new instance of all
func newApp() *app {
	a := app{
		twfUser:  twf.New(),
		twfAdmin: twf.New(),
	}

	// with twf we can change any function for content showing
	a.twfAdmin.MenuFunc = example_templates.AdminMenu

	// if file upload is needed you should use twftemplates.MulipartForm as twf.FormFunc
	a.twfAdmin.FormFunc = twftemplates.MulipartForm

	return &a
}

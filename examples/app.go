package main

import (
	"github.com/tochk/twf"
	"github.com/tochk/twf/examples/example_templates"
)

type app struct {
	twfUser  *twf.TWF
	twfAdmin *twf.TWF
}

func newApp() *app {
	a := app{
		twfUser:  twf.New(),
		twfAdmin: twf.New(),
	}

	a.twfAdmin.MenuFunc = example_templates.AdminMenu
	a.twfAdmin.FormFunc = example_templates.MulipartForm
	return &a
}

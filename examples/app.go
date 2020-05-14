package main

import "github.com/tochk/twf"

type app struct {
	twfUser  *twf.TWF
	twfAdmin *twf.TWF
}

func newApp() *app {
	a := app{
		twfUser:  twf.New(),
		twfAdmin: twf.New(),
	}

	//todo change menu for admin

	return &a
}

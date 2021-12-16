package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tochk/twf"
	"net/http"
	"strconv"
)

type user struct {
	ID      int    `twf:"name:id,title:ID,no_create,no_edit"`
	Login   string `twf:"name:login,title:Login"`
	Email   string `twf:"name:email,title:Email"`
	GroupID int    `twf:"name:group_id,title:Group,fk:0;id;name"`
	Avatar  []byte `twf:"name:avatar,title:Avatar,not_show_on_table,type:file"`
	Edit    string `twf:"title:Edit,name:edit,value:<a href=\"/users/edit/{id}\">Edit</a>,no_create,no_edit,process_parameters"`
}

var users = []user{
	{
		ID:      0,
		Login:   "test1",
		Email:   "me@tochk.ru",
		GroupID: 1,
	},
	{
		ID:      1,
		Login:   "test2",
		Email:   "test2@tochk.ru",
		GroupID: 0,
	},
}

var groups = []group{
	{
		ID:   0,
		Name: "test0",
	},
	{
		ID:   1,
		Name: "test1",
	},
}

type group struct {
	ID   int    `twf:"name:id"`
	Name string `twf:"name:name"`
}

func (a *app) usersListPage(w http.ResponseWriter, r *http.Request) {
	data, err := a.twfAdmin.Table("Users", users, groups)
	if err != nil {
		fmt.Fprint(w, "Err: ", err)
		return
	}
	fmt.Fprint(w, data)
}

func (a *app) usersEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Fprint(w, "Err: ", err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		data, err := a.twfAdmin.EditForm("Users", &users[id], "", groups)
		if err != nil {
			fmt.Fprint(w, "Err: ", err)
			return
		}
		fmt.Fprint(w, data)
	case http.MethodPost:
		r.ParseMultipartForm(32 << 20)
		var user user
		if err := twf.PostFormToStruct(&user, r); err != nil {
			fmt.Fprint(w, "Err: ", err)
			return
		}
		fmt.Println(user)
	}
}

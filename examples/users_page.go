package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tochk/twf"
	"net/http"
	"strconv"
)

type user struct {
	ID    int    `twf:"name:id,title:ID,is_not_creatable,is_not_editable"`
	Login string `twf:"name:login,title:Login"`
	Email string `twf:"name:email,title:Email"`
	//GroupID int    `twf:"name:group_id,title:Group,fk:0,id,name"`
	Avatar []byte `twf:"name:avatar,title:Avatar,is_not_show_on_list,is_not_required,type:file"`
	Edit   string `twf:"process_parameters,title:Edit,name:edit,value:<a href=\"/users/edit/{id}\">Edit</a>,is_not_creatable,is_not_editable"`
}

var users = []user{
	{
		ID:    0,
		Login: "test1",
		Email: "me@tochk.ru",
	},
	{
		ID:    1,
		Login: "test2",
		Email: "test2@tochk.ru",
	},
}

func (a *app) usersListPage(w http.ResponseWriter, r *http.Request) {
	data, err := a.twfAdmin.Table("Users", &user{}, users)
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
		data, err := a.twfAdmin.AddForm("Users", &users[id], "")
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

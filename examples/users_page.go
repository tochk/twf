package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tochk/twf"
	"net/http"
	"strconv"
)

// user structure for user info
type user struct {
	// ID field with disabled creating and editing
	ID int `twf:"name:id,title:ID,no_create,no_edit"`

	Login string `twf:"name:login,title:Login"`
	Email string `twf:"name:email,title:Email"`

	// GroupID field with id of group, will be shown in table as group name because of fk field
	GroupID int `twf:"name:group_id,title:Group,fk:0;id;name"`

	// Avatar field with disabled showing in table
	Avatar []byte `twf:"name:avatar,title:Avatar,not_show_on_table,type:file"`

	// Edit field with disabled creating and editing, but it will be a link in the table
	Edit string `twf:"title:Edit,name:edit,value:<a href=\"/users/edit/{id}\">Edit</a>,no_create,no_edit,process_parameters"`
}

// users - example slice with user info
var users = []user{
	{
		ID:      0,
		Login:   "me",
		Email:   "me@example.com",
		GroupID: 1,
	},
	{
		ID:      1,
		Login:   "not_me",
		Email:   "not_me@example.com",
		GroupID: 0,
	},
}

// usersListPage handler for /users/ page
func (a *app) usersListPage(w http.ResponseWriter, r *http.Request) {
	// print table with values in `users` slice
	// first parameter - page title
	// second parameter - users slice
	// third parameter - groups slice (for linking with GroupID in user struct)
	data, err := a.twfAdmin.Table("Users", users, groups)
	if err != nil {
		fmt.Fprint(w, "Err: ", err)
		return
	}

	// print result to user
	fmt.Fprint(w, data)
}

// usersEditPage handler for /users/edit/{id} page
func (a *app) usersEditPage(w http.ResponseWriter, r *http.Request) {
	// parse id from request
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Fprint(w, "Err: ", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// build add form by loginData structure
		data, err := a.twfAdmin.EditForm("Users", &users[id], "", groups)
		if err != nil {
			fmt.Fprint(w, "Err: ", err)
			return
		}

		// print form page to user
		fmt.Fprint(w, data)
	case http.MethodPost:
		// don't forget to parse form before calling twf.PostFormToStruct function
		r.ParseMultipartForm(32 << 20)

		var user user

		// parse form
		if err := twf.PostFormToStruct(&user, r); err != nil {
			fmt.Fprint(w, "Err: ", err)
			return
		}

		// update user in slice
		users[id] = user

		// redirect to users page
		http.Redirect(w, r, "/users/", http.StatusFound)
	}
}

// groups - example slice with group info
var groups = []group{
	{
		ID:   0,
		Name: "Group with id 0",
	},
	{
		ID:   1,
		Name: "Group with id 1",
	},
}

// group structure for group info
type group struct {
	ID   int    `twf:"name:id"`
	Name string `twf:"name:name"`
}

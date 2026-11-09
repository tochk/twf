// Package twf generates HTML admin pages (tables, add and edit forms) from Go
// structs described by twf struct tags, and parses posted forms back into them.
package twf

import (
	"github.com/tochk/twf/datastruct"
	"github.com/tochk/twf/twftemplates"
)

// Field types with special handling in forms and form parsing.
const (
	typeText     = "text"
	typeCheckbox = "checkbox"
	typeSelect   = "select"
	typeTextarea = "textarea"
	typeFile     = "file"
)

// TWF base struct
type TWF struct {
	HeadFunc         func(string) string
	MenuFunc         func() string
	ListFunc         func([]datastruct.Field, string) string
	ListItemFunc     func([]any) string
	FooterFunc       func() string
	FormFunc         func(string, string) string
	FormItemText     func(datastruct.Field) string
	FormItemTextarea func(datastruct.Field) string
	FormItemSelect   func(datastruct.Field, []datastruct.FkKV, any) string
	FormItemCheckbox func(datastruct.Field) string
}

// New creates new TWF instance with default parameters
func New() *TWF {
	twf := TWF{
		HeadFunc:         twftemplates.Head,
		MenuFunc:         twftemplates.Menu,
		ListFunc:         twftemplates.ItemList,
		ListItemFunc:     twftemplates.Item,
		FooterFunc:       twftemplates.Footer,
		FormFunc:         twftemplates.Form,
		FormItemText:     twftemplates.FormItem,
		FormItemTextarea: twftemplates.FormItemTextarea,
		FormItemSelect:   twftemplates.FormItemSelect,
		FormItemCheckbox: twftemplates.FormItemCheckbox,
	}
	return &twf
}

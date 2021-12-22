package twf

import (
	"github.com/tochk/twf/datastruct"
	"github.com/tochk/twf/twftemplates"
)

// TWF base struct
type TWF struct {
	HeadFunc         func(string) string
	MenuFunc         func() string
	ListFunc         func([]datastruct.Field, string) string
	ListItemFunc     func([]interface{}) string
	FooterFunc       func() string
	FormFunc         func(string, string) string
	FormItemText     func(datastruct.Field) string
	FormItemTextarea func(datastruct.Field) string
	FormItemSelect   func(datastruct.Field, []datastruct.FkKV, interface{}) string
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

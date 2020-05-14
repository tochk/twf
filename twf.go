package twf

import (
	"github.com/tochk/twf/datastruct"
	"github.com/tochk/twf/twftemplates"
)

type TWF struct {
	HeadFunc         func(string) string
	MenuFunc         func() string
	ListFunc         func([]datastruct.Field, string) string
	ListItemFunc     func([]interface{}) string
	FooterFunc       func() string
	FormFunc         func(string, string) string
	FormItemFunc     func(datastruct.Field) string
	FormItemSelect   func(field datastruct.Field, kvs []datastruct.FkKV, selectedID interface{}) string
	FormItemCheckbox func(field datastruct.Field) string
}

func New() *TWF {
	twf := TWF{
		HeadFunc:         twftemplates.Head,
		MenuFunc:         twftemplates.Menu,
		ListFunc:         twftemplates.ItemList,
		ListItemFunc:     twftemplates.Item,
		FooterFunc:       twftemplates.Footer,
		FormFunc:         twftemplates.Form,
		FormItemFunc:     twftemplates.FormItem,
		FormItemSelect:   twftemplates.FormItemSelect,
		FormItemCheckbox: twftemplates.FormItemCheckbox,
	}
	return &twf
}

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
	FormItemText     func(datastruct.Field) string
	FormItemTextarea func(datastruct.Field) string
	FormItemSelect   func(datastruct.Field, []datastruct.FkKV, interface{}) string
	FormItemCheckbox func(datastruct.Field) string
}

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

func filterTableFields(fields []datastruct.Field) []datastruct.Field {
	res := make([]datastruct.Field, 0, len(fields))
	for _, field := range fields {
		if field.NotShowOnTable {
			continue
		}
		res = append(res, field)
	}

	return res
}

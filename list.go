package twf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tochk/twf/datastruct"
	"github.com/tochk/twf/twftemplates"
)

var (
	HeadFunc         func(string) string                                                                = twftemplates.Head
	MenuFunc         func() string                                                                      = twftemplates.Menu
	ListFunc         func([]datastruct.Field, string) string                                            = twftemplates.ItemList
	ListItemFunc     func([]interface{}) string                                                         = twftemplates.Item
	FooterFunc       func() string                                                                      = twftemplates.Footer
	FormFunc         func(string, string) string                                                        = twftemplates.Form
	FormItemFunc     func(datastruct.Field) string                                                      = twftemplates.FormItem
	FormItemSelect   func(field datastruct.Field, kvs []datastruct.FkKV, selectedID interface{}) string = twftemplates.FormItemSelect
	FormItemCheckbox func(field datastruct.Field) string                                                = twftemplates.FormItemCheckbox
)

var (
	errorNotSlice        = errors.New("items must be slice")
	ErrFksIndexNotExists = errors.New("fks index not exists")
	ErrFksMustBeSLice    = errors.New("fks must me a slice")
)

func processParameters(value string, fields map[string]string) interface{} {
	for k, v := range fields {
		value = strings.Replace(value, "{"+k+"}", v, -1)
	}
	return value
}

func List(title string, isAdmin bool, item interface{}, items interface{}, fks ...interface{}) (string, error) {
	fields, err := GetFieldDescription(item)
	if err != nil {
		return "", err
	}
	res := strings.Builder{}
	res.WriteString(HeadFunc(title))
	res.WriteString(MenuFunc())
	content := strings.Builder{}
	switch reflect.TypeOf(items).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(items)
		for i := 0; i < s.Len(); i++ {
			itemsSlice := make([]interface{}, 0, s.Index(i).NumField())
			data := map[string]string{}
			for j := 0; j < s.Index(i).NumField(); j++ {
				var value interface{}
				if fields[j].Value == "" {
					tmp := s.Index(i).Field(j)
					if s.Index(i).Field(j).Kind() == reflect.Ptr {
						if tmp.IsNil() {
							value = "nil"
						} else {
							value = tmp.Elem().Interface()
						}
					} else {
						value = tmp.Interface()
					}
				} else {
					if fields[j].ProcessParameters {
						value = processParameters(fields[j].Value, data)
					} else {
						value = fields[j].Value
					}
				}
				var fkValue interface{}
				if fields[j].FkInfo != nil {
					fksInfo := fields[j].FkInfo
					if len(fks) <= fksInfo.FksIndex {
						return "", ErrFksIndexNotExists
					}
					fksSlice := fks[fksInfo.FksIndex]
					if reflect.TypeOf(fksSlice).Kind() != reflect.Slice {
						return "", ErrFksMustBeSLice
					}
					fksValue := reflect.ValueOf(fksSlice)
					for k := 0; k < fksValue.Len(); k++ {
						v := fksValue.Index(k)
						fkKv := datastruct.FkKV{}
						for l := 0; l < v.NumField(); l++ {
							fv := v.Field(l)
							ft := v.Type().Field(l)
							if tag, ok := ft.Tag.Lookup("twf"); ok {
								tags := strings.Split(tag, ",")
								for _, e := range tags {
									if len(e) > 5 {
										if e[:5] == "name:" {
											if e[5:] == fields[j].FkInfo.Name {
												fkKv.Name = fv
											}
											if e[5:] == fields[j].FkInfo.ID {
												fkKv.ID = fv.Interface()
											}
										}
									}
								}
							} else {
								return "", ErrParameterNotFound
							}
						}
						if value == fkKv.ID {
							fkValue = value
							value = fkKv.Name
							break
						}
					}
				}
				if fkValue != nil {
					data[fields[j].Name] = fmt.Sprint(fkValue)
				} else {
					data[fields[j].Name] = fmt.Sprint(value)
				}
				if !fields[j].IsNotShowOnList {
					itemsSlice = append(itemsSlice, value)
				}
			}
			content.WriteString(ListItemFunc(itemsSlice))
		}
	default:
		return "", errorNotSlice
	}
	f2 := make([]datastruct.Field, 0, len(fields))
	for _, e := range fields {
		if e.IsNotShowOnList {
			continue
		}
		f2 = append(f2, e)
	}
	fields = f2
	res.WriteString(ListFunc(fields, content.String()))
	res.WriteString(FooterFunc())
	return res.String(), nil
}

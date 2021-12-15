package twf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tochk/twf/datastruct"
)

func processParameters(value string, fields map[string]string) interface{} {
	for k, v := range fields {
		value = strings.Replace(value, "{"+k+"}", v, -1)
	}
	return value
}

func (t *TWF) Table(title string, item interface{}, items interface{}, fks ...interface{}) (string, error) {
	fields, err := GetFieldDescription(item)
	if err != nil {
		return "", err
	}
	res := strings.Builder{}
	res.WriteString(t.HeadFunc(title))
	res.WriteString(t.MenuFunc())
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
				if !fields[j].NotShowOnTable {
					itemsSlice = append(itemsSlice, value)
				}
			}
			content.WriteString(t.ListItemFunc(itemsSlice))
		}
	default:
		return "", errorNotSlice
	}
	f2 := make([]datastruct.Field, 0, len(fields))
	for _, e := range fields {
		if e.NotShowOnTable {
			continue
		}
		f2 = append(f2, e)
	}
	fields = f2
	res.WriteString(t.ListFunc(fields, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

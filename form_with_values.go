package twf

import (
	"fmt"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"strings"
)

func (t *TWF) FormWithValues(title string, item interface{}, link string, fks ...interface{}) (string, error) {
	fields, err := GetFieldDescription(item)
	if err != nil {
		return "", err
	}
	res := strings.Builder{}
	res.WriteString(t.HeadFunc(title))
	res.WriteString(t.MenuFunc())
	content := strings.Builder{}
	switch reflect.TypeOf(item).Kind() {
	case reflect.Ptr:
		s := reflect.ValueOf(item).Elem()
		data := map[string]string{}
		for i := 0; i < s.NumField(); i++ {
			var value interface{}
			field := fields[i]

			if fields[i].Value == "" {
				tmp := s.Field(i)
				if tmp.Kind() == reflect.Ptr {
					if tmp.IsNil() {
						value = ""
					} else {
						value = tmp.Elem().Interface()
					}
				} else {
					value = tmp.Interface()
				}
			} else {
				if fields[i].ProcessParameters {
					value = processParameters(fields[i].Value, data)
				} else {
					value = fields[i].Value
				}
			}

			kvs := make([]datastruct.FkKV, 0)
			if fields[i].FkInfo != nil {
				fksInfo := fields[i].FkInfo
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
										if e[5:] == fields[i].FkInfo.Name {
											fkKv.Name = fv
										}
										if e[5:] == fields[i].FkInfo.ID {
											fkKv.ID = fv.Interface()
										}
									}
								}
							}
						} else {
							return "", ErrParameterNotFound
						}
					}
					kvs = append(kvs, fkKv)
				}
			}

			field.Value = fmt.Sprint(value)
			if !field.IsNotEditable {
				switch field.Type {
				case "select":
					content.WriteString(t.FormItemSelect(field, kvs, value))
				case "checkbox":
					content.WriteString(t.FormItemCheckbox(field))
				case "textarea":
					content.WriteString(t.FormItemTextarea(field))
				default:
					content.WriteString(t.FormItemText(field))
				}
			}
		}
	default:
		return "", errorNotStruct
	}
	res.WriteString(t.FormFunc(link, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

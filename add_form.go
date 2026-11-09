package twf

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// AddForm returns page with add form. Values from `value` tag are used as defaults, otherwise values from item are used
func (t *TWF) AddForm(title string, item interface{}, link string, fks ...interface{}) (string, error) {
	s, err := structPtrValue("twf.AddForm", item)
	if err != nil {
		return "", err
	}

	if err := validateFks(fks); err != nil {
		return "", errors.Wrap(err, "twf.AddForm")
	}

	fields, err := getFieldDescription(reflect.TypeOf(item))
	if err != nil {
		return "", errors.Wrap(err, "twf.AddForm")
	}
	res := strings.Builder{}
	res.WriteString(t.HeadFunc(title))
	res.WriteString(t.MenuFunc())
	content := strings.Builder{}

	data := map[string]string{}
	for i := 0; i < s.NumField(); i++ {
		var value interface{}
		field := fields[i]
		switch {
		case field.Value == "":
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
		case field.ProcessParameters:
			value = processParameters(field.Value, data)
		default:
			value = field.Value
		}

		kvs, err := getFKSlice(field.FkInfo, fks...)
		if err != nil {
			return "", errors.Wrap(err, "twf.AddForm")
		}

		field.Value = fmt.Sprint(value)
		data[field.Name] = field.Value

		if field.NoCreate {
			continue
		}

		switch field.Type {
		case "select":
			content.WriteString(t.FormItemSelect(field, kvs, nil))
		case "checkbox":
			content.WriteString(t.FormItemCheckbox(field))
		case "textarea":
			content.WriteString(t.FormItemTextarea(field))
		default:
			content.WriteString(t.FormItemText(field))
		}
	}

	res.WriteString(t.FormFunc(link, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

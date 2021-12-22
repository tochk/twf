package twf

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// EditForm returns page with add form (without values, if you want a form with values from struct - use AddForm)
func (t *TWF) EditForm(pageTitle string, item interface{}, link string, fks ...interface{}) (string, error) {
	itemType := reflect.TypeOf(item)
	if itemType.Kind() != reflect.Ptr {
		return "", fmt.Errorf("twf.getSliceElementPtrType: expected ptr to struct, got %s", itemType.Kind().String())
	}
	if itemType.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("twf.getSliceElementPtrType: expected ptr to struct, got ptr to %s", itemType.Elem().Kind().String())
	}

	fields, err := getFieldDescription(reflect.TypeOf(item))
	if err != nil {
		return "", errors.Wrap(err, "twf.EditForm")
	}
	res := strings.Builder{}
	res.WriteString(t.HeadFunc(pageTitle))
	res.WriteString(t.MenuFunc())
	content := strings.Builder{}

	s := reflect.ValueOf(item).Elem()
	data := map[string]string{}
	for i := 0; i < s.NumField(); i++ {
		var value interface{}
		field := fields[i]

		if fields[i].Value == "" {
			value = getFieldValue(s.Field(i))
		} else {
			if fields[i].ProcessParameters {
				value = processParameters(fields[i].Value, data)
			} else {
				value = fields[i].Value
			}
		}

		kvs, err := getFKSlice(fields[i].FkInfo, fks...)
		if err != nil {
			return "", errors.Wrap(err, "twf.EditForm")
		}

		field.Value = fmt.Sprint(value)

		if field.NoEdit {
			continue
		}

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

	res.WriteString(t.FormFunc(link, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

package twf

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// EditForm returns page with edit form filled with values from item. `value` tag is used only for fields with
// process_parameters, real values are never overridden
func (t *TWF) EditForm(pageTitle string, item interface{}, link string, fks ...interface{}) (string, error) {
	s, err := structPtrValue("twf.EditForm", item)
	if err != nil {
		return "", err
	}

	if err := validateFks(fks); err != nil {
		return "", errors.Wrap(err, "twf.EditForm")
	}

	fields, err := getFieldDescription(reflect.TypeOf(item))
	if err != nil {
		return "", errors.Wrap(err, "twf.EditForm")
	}
	res := strings.Builder{}
	res.WriteString(t.HeadFunc(pageTitle))
	res.WriteString(t.MenuFunc())
	content := strings.Builder{}

	data := map[string]string{}
	for i := 0; i < s.NumField(); i++ {
		var value interface{}
		field := fields[i]

		if field.ProcessParameters && field.Value != "" {
			value = processParameters(field.Value, data)
		} else {
			value = getFieldValue(s.Field(i))
		}

		kvs, err := getFKSlice(field.FkInfo, fks...)
		if err != nil {
			return "", errors.Wrap(err, "twf.EditForm")
		}

		field.Value = fmt.Sprint(value)
		data[field.Name] = field.Value

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

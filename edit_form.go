package twf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// EditForm returns page with edit form filled with values from item. `value` tag is used only for fields with
// process_parameters, real values are never overridden
func (t *TWF) EditForm(pageTitle string, item any, link string, fks ...any) (string, error) {
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
		var value any
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
		case typeSelect:
			content.WriteString(t.FormItemSelect(field, kvs, value))
		case typeCheckbox:
			content.WriteString(t.FormItemCheckbox(field))
		case typeTextarea:
			content.WriteString(t.FormItemTextarea(field))
		default:
			content.WriteString(t.FormItemText(field))
		}
	}

	res.WriteString(t.FormFunc(link, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

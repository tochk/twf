package twf

import (
	"fmt"
	"reflect"
	"strings"
)

func processParameters(value string, fields map[string]string) interface{} {
	for k, v := range fields {
		value = strings.Replace(value, "{"+k+"}", v, -1)
	}
	return value
}

func (t *TWF) Table(title string, items interface{}, fks ...interface{}) (string, error) {
	if reflect.TypeOf(items).Kind() != reflect.Slice {
		return "", fmt.Errorf("twf.Table: expected slice, got %s", reflect.TypeOf(items).Kind().String())
	}

	for i, fkSlice := range fks {
		if reflect.TypeOf(fkSlice).Kind() != reflect.Slice {
			return "", fmt.Errorf("twf.Table: error on fks idx %d expected slice, got %s", i, reflect.TypeOf(fkSlice).Kind().String())
		}
	}

	item, err := getSliceElementPtrType(items)
	if err != nil {
		return "", err
	}

	fields, err := getFieldDescription(item)
	if err != nil {
		return "", err
	}

	res := strings.Builder{}
	res.WriteString(t.HeadFunc(title))
	res.WriteString(t.MenuFunc())
	content := strings.Builder{}

	s := reflect.ValueOf(items)
	for i := 0; i < s.Len(); i++ {
		itemsSlice := make([]interface{}, 0, s.Index(i).NumField())
		data := map[string]string{}
		for j := 0; j < s.Index(i).NumField(); j++ {
			var value interface{}
			if fields[j].Value == "" {
				tmp := s.Index(i).Field(j)
				value = getFieldValue(tmp)
			} else {
				if fields[j].ProcessParameters {
					value = processParameters(fields[j].Value, data)
				} else {
					value = fields[j].Value
				}
			}

			var fkValue interface{}
			fkValue, value, err = getFKValue(fields[j].FkInfo, value, fks...)
			if err != nil {
				return "", err
			}

			if fkValue != nil {
				data[fields[j].Name] = fmt.Sprint(fkValue)
				if !fields[j].NotShowOnTable {
					itemsSlice = append(itemsSlice, fkValue)
				}
				continue
			}

			data[fields[j].Name] = fmt.Sprint(value)
			if !fields[j].NotShowOnTable {
				itemsSlice = append(itemsSlice, value)
			}
		}
		content.WriteString(t.ListItemFunc(itemsSlice))
	}

	fields = filterTableFields(fields)

	res.WriteString(t.ListFunc(fields, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

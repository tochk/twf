package twf

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"strings"
)

// processParameters replaces {field} in value
func processParameters(value string, fields map[string]string) interface{} {
	for k, v := range fields {
		value = strings.Replace(value, "{"+k+"}", v, -1)
	}
	return value
}

// filterTableFields filters fields is TWF.NotShowOnTable flag is enabled
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

// Table return table with contents of given slice
func (t *TWF) Table(title string, slice interface{}, fks ...interface{}) (string, error) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return "", fmt.Errorf("twf.Table: expected slice, got %s", reflect.TypeOf(slice).Kind().String())
	}

	for i, fkSlice := range fks {
		if reflect.TypeOf(fkSlice).Kind() != reflect.Slice {
			return "", fmt.Errorf("twf.Table: error on fks idx %d expected slice, got %s", i, reflect.TypeOf(fkSlice).Kind().String())
		}
	}

	item, err := getSliceElementPtrType(slice)
	if err != nil {
		return "", errors.Wrap(err, "twf.Table")
	}

	fields, err := getFieldDescription(item)
	if err != nil {
		return "", errors.Wrap(err, "twf.Table")
	}

	res := strings.Builder{}
	res.WriteString(t.HeadFunc(title))
	res.WriteString(t.MenuFunc())
	content := strings.Builder{}

	s := reflect.ValueOf(slice)
	for i := 0; i < s.Len(); i++ {
		itemsSlice, err := generateTableItemsSlice(s, fields, i, fks)
		if err != nil {
			return "", err
		}
		content.WriteString(t.ListItemFunc(itemsSlice))
	}

	fields = filterTableFields(fields)

	res.WriteString(t.ListFunc(fields, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

func generateTableItemsSlice(s reflect.Value, fields []datastruct.Field, i int, fks []interface{}) ([]interface{}, error) {
	itemsSlice := make([]interface{}, 0, s.Index(i).NumField())
	data := map[string]string{}
	for j := 0; j < s.Index(i).NumField(); j++ {
		var value interface{}

		if fields[j].ProcessParameters {
			value = processParameters(fields[j].Value, data)
		} else {
			tmp := s.Index(i).Field(j)
			value = getFieldValue(tmp)
		}

		var fkValue interface{}
		fkValue, value, err := getFKValue(fields[j].FkInfo, value, fks)
		if err != nil {
			return nil, err
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

	return itemsSlice, nil
}

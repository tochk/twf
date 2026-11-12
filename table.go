package twf

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/tochk/twf/datastruct"
)

// processParameters replaces {field} in value. Replacement is done in a single pass,
// so substituted values are never expanded again
func processParameters(value string, fields map[string]string) string {
	if len(fields) == 0 {
		return value
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	oldNew := make([]string, 0, len(keys)*2)
	for _, k := range keys {
		oldNew = append(oldNew, "{"+k+"}", fields[k])
	}
	return strings.NewReplacer(oldNew...).Replace(value)
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
func (t *TWF) Table(title string, slice any, fks ...any) (string, error) {
	item, err := getSliceElementPtrType(slice)
	if err != nil {
		return "", fmt.Errorf("twf.Table: %w", err)
	}

	if err := validateFks(fks); err != nil {
		return "", fmt.Errorf("twf.Table: %w", err)
	}

	fields, err := getFieldDescription(item)
	if err != nil {
		return "", fmt.Errorf("twf.Table: %w", err)
	}

	res := strings.Builder{}
	res.WriteString(t.HeadFunc(title))
	res.WriteString(t.MenuFunc())
	content := strings.Builder{}

	s := reflect.ValueOf(slice)
	for i := 0; i < s.Len(); i++ {
		itemsSlice, err := generateTableItemsSlice(s, fields, i, fks)
		if err != nil {
			return "", fmt.Errorf("twf.Table: %w", err)
		}
		content.WriteString(t.ListItemFunc(itemsSlice))
	}

	fields = filterTableFields(fields)

	res.WriteString(t.ListFunc(fields, content.String()))
	res.WriteString(t.FooterFunc())
	return res.String(), nil
}

func generateTableItemsSlice(s reflect.Value, fields []datastruct.Field, i int, fks []any) ([]any, error) {
	itemsSlice := make([]any, 0, s.Index(i).NumField())
	data := map[string]string{}
	for j := 0; j < s.Index(i).NumField(); j++ {
		var value any

		if fields[j].ProcessParameters {
			value = processParameters(fields[j].Value, data)
		} else {
			tmp := s.Index(i).Field(j)
			value = getFieldValue(tmp)
		}

		var fkValue any
		fkValue, value, err := getFKValue(fields[j].FkInfo, value, fks)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldLabel(s.Type().Elem(), s.Type().Elem().Field(j).Name), err)
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

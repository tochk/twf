package twf

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"strconv"
	"strings"
)

// getSliceElementPtrType returns pointer to slice element
func getSliceElementPtrType(slice interface{}) (reflect.Type, error) {
	s := reflect.TypeOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("twf.getSliceElementPtrType: expected slice, got %s", s.Kind().String())
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, fmt.Errorf("twf.getSliceElementPtrType: slice containts %s, expected struct", s.Kind().String())
	}
	s = reflect.PtrTo(s)
	return s, nil
}

// getFieldDescription returns information of all fields in given struct
func getFieldDescription(s reflect.Type) ([]datastruct.Field, error) {
	if s.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("twf.getFieldDescription: expected ptr to struct, got %s", s.Kind().String())
	}

	if s.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("twf.getFieldDescription: expected ptr to struct, got ptr to %s", s.Elem().Kind().String())
	}

	s = s.Elem()

	fields := make([]datastruct.Field, 0, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		field := datastruct.Field{
			Type: "text",
		}
		f := s.Field(i)
		switch f.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.Type = "number"
		case reflect.Bool:
			field.Type = "checkbox"
		}
		tagContent := strings.Split(f.Tag.Get("twf"), ",")
		for _, e := range tagContent {
			switch {
			case strings.HasPrefix(e, "fk:"):
				fkInfo := strings.Split(e[3:], ";")
				if len(fkInfo) != 3 {
					return nil, fmt.Errorf("twf.getFieldDescription: invalid fk info {%s}, fk must be {fk_slice_index;id;name}", e[3:])
				}
				fkID, err := strconv.Atoi(fkInfo[0])
				if err != nil {
					return nil, errors.Wrap(err, "twf.getFieldDescription: can't convert fk id to int")
				}
				field.FkInfo = &datastruct.FkInfo{
					FksIndex: fkID,
					ID:       fkInfo[1],
					Name:     fkInfo[2],
				}
				field.Type = "select"
			case strings.HasPrefix(e, "name:"):
				field.Name = e[5:]
			case strings.HasPrefix(e, "type:"):
				field.Type = e[5:]
			case strings.HasPrefix(e, "title:"):
				field.Title = e[6:]
			case strings.HasPrefix(e, "value:"):
				field.Value = e[6:]
			case strings.HasPrefix(e, "placeholder:"):
				field.Placeholder = e[12:]
			case e == "no_create":
				field.NoCreate = true
			case e == "no_edit":
				field.NoEdit = true
			case e == "required":
				field.Required = true
			case e == "not_show_on_table":
				field.NotShowOnTable = true
			case e == "process_parameters":
				field.ProcessParameters = true
			default:
				return nil, fmt.Errorf("twf.getFieldDescription: unexpected parameter in twf tag: %s", e)
			}
		}
		fields = append(fields, field)
	}
	return fields, nil
}

// getFieldValue returns value of the field. Returns <nil> if value is nil
func getFieldValue(field reflect.Value) interface{} {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return "<nil>"
		} else {
			return field.Elem().Interface()
		}
	}

	return field.Interface()
}

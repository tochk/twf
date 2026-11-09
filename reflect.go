package twf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/tochk/twf/datastruct"
)

// getSliceElementPtrType returns pointer to slice element
func getSliceElementPtrType(slice any) (reflect.Type, error) {
	s := reflect.TypeOf(slice)
	if s == nil {
		return nil, errors.New("twf.getSliceElementPtrType: expected slice, got nil")
	}
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("twf.getSliceElementPtrType: expected slice, got %s", s.Kind().String())
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, fmt.Errorf("twf.getSliceElementPtrType: slice contains %s, expected struct", s.Kind().String())
	}
	s = reflect.PointerTo(s)
	return s, nil
}

// getFieldDescription returns information of all fields in given struct
func getFieldDescription(s reflect.Type) ([]datastruct.Field, error) {
	if s.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("twf.getFieldDescription: expected ptr to struct, got %s", s.Kind().String())
	}

	if s.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("twf.getFieldDescription: expected ptr to struct, got ptr to %s", s.Elem().Kind().String())
	}

	s = s.Elem()

	fields := make([]datastruct.Field, 0, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		field := datastruct.Field{
			Type: typeText,
		}
		f := s.Field(i)
		if !f.IsExported() {
			return nil, fmt.Errorf("twf.getFieldDescription: field %s must be exported", f.Name)
		}
		switch f.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.Type = "number"
		case reflect.Bool:
			field.Type = typeCheckbox
		}
		if t, ok := f.Tag.Lookup("twf"); !ok || t == "" {
			return nil, errors.New("twf tag must be present and not be empty")
		}
		tagContent := strings.SplitSeq(f.Tag.Get("twf"), ",")
		for e := range tagContent {
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
				if fkID < 0 {
					return nil, fmt.Errorf("twf.getFieldDescription: fk index must not be negative, got %d", fkID)
				}
				fkType := f.Type
				if fkType.Kind() == reflect.Pointer {
					fkType = fkType.Elem()
				}
				if !fkType.Comparable() {
					return nil, fmt.Errorf("twf.getFieldDescription: fk field %s must be of comparable type, got %s", f.Name, f.Type.String())
				}
				field.FkInfo = &datastruct.FkInfo{
					FksIndex: fkID,
					ID:       fkInfo[1],
					Name:     fkInfo[2],
				}
				field.Type = typeSelect
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
			case e == "disabled":
				field.Disabled = true
			default:
				return nil, fmt.Errorf("twf.getFieldDescription: unexpected parameter in twf tag: %q", e)
			}
		}
		fields = append(fields, field)
	}
	return fields, nil
}

// structPtrValue checks that item is a non-nil pointer to struct and returns the struct value
func structPtrValue(fn string, item any) (reflect.Value, error) {
	t := reflect.TypeOf(item)
	if t == nil {
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got nil", fn)
	}
	if t.Kind() != reflect.Pointer {
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got %s", fn, t.Kind().String())
	}
	if t.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got ptr to %s", fn, t.Elem().Kind().String())
	}
	v := reflect.ValueOf(item)
	if v.IsNil() {
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got nil ptr", fn)
	}
	return v.Elem(), nil
}

// getFieldValue returns value of the field. Returns <nil> if value is nil
func getFieldValue(field reflect.Value) any {
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return "<nil>"
		}
		return field.Elem().Interface()
	}

	return field.Interface()
}

package twf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tochk/twf/datastruct"
)

// getSliceElementPtrType returns pointer to slice element
func getSliceElementPtrType(slice any) (reflect.Type, error) {
	s := reflect.TypeOf(slice)
	if s == nil {
		return nil, errors.New("twf.getSliceElementPtrType: expected slice of structs, got nil")
	}
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("twf.getSliceElementPtrType: expected slice of structs, got %s", s.Kind().String())
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		hint := ""
		if s.Kind() == reflect.Pointer && s.Elem().Kind() == reflect.Struct {
			hint = " (use []T, not []*T)"
		}
		return nil, fmt.Errorf("twf.getSliceElementPtrType: expected slice of structs, got slice of %s%s", s.Kind().String(), hint)
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
	for f := range s.Fields() {
		field, err := parseField(&f)
		if err != nil {
			return nil, fmt.Errorf("twf.getFieldDescription: field %s: %w", fieldLabel(s, f.Name), err)
		}
		fields = append(fields, field)
	}
	return fields, nil
}

// fieldLabel returns Type.Field, or just Field for anonymous struct types
func fieldLabel(t reflect.Type, name string) string {
	if t.Name() == "" {
		return name
	}
	return t.Name() + "." + name
}

// parseField builds field description from struct field and its twf tag
func parseField(f *reflect.StructField) (datastruct.Field, error) {
	field := datastruct.Field{
		Type: typeText,
	}
	if !f.IsExported() {
		return field, errors.New("must be exported")
	}
	switch f.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.Type = "number"
	case reflect.Float32, reflect.Float64:
		field.Type = "number"
		field.Step = "any"
	case reflect.Bool:
		field.Type = typeCheckbox
	}
	if t, ok := f.Tag.Lookup("twf"); !ok || t == "" {
		return field, errors.New("twf tag missing or empty")
	}
	for e := range strings.SplitSeq(f.Tag.Get("twf"), ",") {
		switch {
		case strings.HasPrefix(e, "fk:"):
			fkInfo := strings.Split(e[3:], ";")
			if len(fkInfo) != 3 {
				return field, fmt.Errorf("invalid fk %q, expected fk:<index>;<idField>;<nameField>", e)
			}
			fkID, err := strconv.Atoi(fkInfo[0])
			if err != nil {
				return field, fmt.Errorf("fk index %q is not an integer: %w", fkInfo[0], err)
			}
			if fkID < 0 {
				return field, fmt.Errorf("fk index must not be negative, got %d", fkID)
			}
			fkType := f.Type
			if fkType.Kind() == reflect.Pointer {
				fkType = fkType.Elem()
			}
			if !fkType.Comparable() {
				return field, fmt.Errorf("fk field must be comparable, got %s", f.Type.String())
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
			return field, fmt.Errorf("unknown twf tag option %q", e)
		}
	}
	return field, nil
}

// structPtrValue checks that item is a non-nil pointer to struct and returns the struct value
func structPtrValue(fn string, item any) (reflect.Value, error) {
	t := reflect.TypeOf(item)
	if t == nil {
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got nil", fn)
	}
	if t.Kind() != reflect.Pointer {
		hint := ""
		if t.Kind() == reflect.Struct {
			hint = " (pass &v)"
		}
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got %s%s", fn, t.Kind().String(), hint)
	}
	if t.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got ptr to %s", fn, t.Elem().Kind().String())
	}
	v := reflect.ValueOf(item)
	if v.IsNil() {
		return reflect.Value{}, fmt.Errorf("%s: expected ptr to struct, got nil %s", fn, t.String())
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

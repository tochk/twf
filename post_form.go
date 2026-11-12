package twf

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/tochk/twf/datastruct"
)

// absentMeansFalse reports whether a bool field missing from the form must be set to false.
// Browsers don't submit unchecked checkboxes, but they don't submit disabled inputs either,
// and no_edit fields are not rendered in the edit form at all, so those are left unchanged
func absentMeansFalse(field *datastruct.Field, kind reflect.Kind) bool {
	return kind == reflect.Bool && field.Type == typeCheckbox && !field.Disabled && !field.NoEdit
}

// PostFormToStruct marshal post form to given struct
func PostFormToStruct(item any, r *http.Request) error {
	s, err := structPtrValue("twf.PostFormToStruct", item)
	if err != nil {
		return err
	}

	if r == nil {
		return errors.New("twf.PostFormToStruct: request is nil")
	}

	fields, err := getFieldDescription(reflect.TypeOf(item))
	if err != nil {
		return fmt.Errorf("twf.PostFormToStruct: %w", err)
	}
	// unparsed form means no data at all, nothing is considered absent
	formParsed := r.PostForm != nil
	postFormMap := map[string]string{}
	for k := range r.PostForm {
		postFormMap[k] = r.PostForm.Get(k)
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		formValue, ok := postFormMap[fields[i].Name]
		if !ok && fields[i].Type != typeFile && (!formParsed || !absentMeansFalse(&fields[i], f.Kind())) {
			continue
		}
		// the empty option of a select means "nothing selected"
		if fields[i].Type == typeSelect && formValue == "" {
			f.SetZero()
			continue
		}
		switch f.Kind() {
		case reflect.String:
			f.SetString(formValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(formValue, 10, f.Type().Bits())
			if err != nil {
				return fmt.Errorf("twf.PostFormToStruct: %w", err)
			}
			f.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(formValue, 10, f.Type().Bits())
			if err != nil {
				return fmt.Errorf("twf.PostFormToStruct: %w", err)
			}
			f.SetUint(v)
		case reflect.Float64, reflect.Float32:
			v, err := strconv.ParseFloat(formValue, f.Type().Bits())
			if err != nil {
				return fmt.Errorf("twf.PostFormToStruct: %w", err)
			}
			f.SetFloat(v)
		case reflect.Bool:
			f.SetBool(formValue == "on")
		case reflect.Slice:
			if f.Type() == reflect.TypeFor[[]byte]() && fields[i].Type == typeFile {
				file, _, err := r.FormFile(fields[i].Name)
				if errors.Is(err, http.ErrMissingFile) {
					continue
				}
				if err != nil {
					return fmt.Errorf("twf.PostFormToStruct: %w", err)
				}
				data, err := io.ReadAll(file)
				if err != nil {
					return fmt.Errorf("twf.PostFormToStruct: %w", err)
				}
				f.SetBytes(data)
			}
		}
	}
	return nil
}

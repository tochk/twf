package twf

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"strings"
)

// validateFk checks that fk is a slice of structs whose fields are all exported and have a twf tag
func validateFk(idx int, fk interface{}) error {
	t := reflect.TypeOf(fk)
	if t == nil {
		return fmt.Errorf("twf.validateFk: fks with index %d must be slice, got nil", idx)
	}
	if t.Kind() != reflect.Slice {
		return fmt.Errorf("twf.validateFk: fks with index %d must be slice, not %s", idx, t.Kind().String())
	}
	e := t.Elem()
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("twf.validateFk: fks with index %d must be slice of structs, got slice of %s", idx, e.Kind().String())
	}
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		if !f.IsExported() {
			return fmt.Errorf("twf.validateFk: fks with index %d: field %s must be exported", idx, f.Name)
		}
		if _, ok := f.Tag.Lookup("twf"); !ok {
			return fmt.Errorf("twf.validateFk: fks with index %d: twf tag with twf:\"name:value\" must be present in all fks structs, missing on field %s", idx, f.Name)
		}
	}
	return nil
}

// validateFks checks every element of fks with validateFk
func validateFks(fks []interface{}) error {
	for i, fk := range fks {
		if err := validateFk(i, fk); err != nil {
			return err
		}
	}
	return nil
}

// fkSliceValue returns validated fks slice selected by fksInfo
func fkSliceValue(fksInfo *datastruct.FkInfo, fks []interface{}) (reflect.Value, error) {
	if fksInfo.FksIndex < 0 || len(fks) <= fksInfo.FksIndex {
		return reflect.Value{}, fmt.Errorf("fks index %d not provided", fksInfo.FksIndex)
	}
	fk := fks[fksInfo.FksIndex]
	if err := validateFk(fksInfo.FksIndex, fk); err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(fk), nil
}

// getFKKV extracts id and name of fks element v. Fields are matched by name tag
func getFKKV(v reflect.Value, fksInfo *datastruct.FkInfo) datastruct.FkKV {
	fkKv := datastruct.FkKV{}
	for l := 0; l < v.NumField(); l++ {
		for _, e := range strings.Split(v.Type().Field(l).Tag.Get("twf"), ",") {
			if !strings.HasPrefix(e, "name:") {
				continue
			}
			if e[5:] == fksInfo.Name {
				fkKv.Name = v.Field(l)
			}
			if e[5:] == fksInfo.ID {
				fkKv.ID = v.Field(l).Interface()
			}
		}
	}
	return fkKv
}

// getFKValue returns fk value and original value
func getFKValue(fksInfo *datastruct.FkInfo, originalValue interface{},
	fks []interface{}) (fkValue interface{}, value interface{}, err error) {

	if fksInfo == nil {
		return nil, originalValue, nil
	}

	fksValue, err := fkSliceValue(fksInfo, fks)
	if err != nil {
		return nil, nil, errors.Wrap(err, "twf.getFKValue")
	}

	for k := 0; k < fksValue.Len(); k++ {
		fkKv := getFKKV(fksValue.Index(k), fksInfo)
		if originalValue == fkKv.ID {
			return fkKv.Name, originalValue, nil
		}
	}

	return nil, originalValue, nil
}

// getFKSlice return all fks of given fks info
func getFKSlice(fksInfo *datastruct.FkInfo, fks ...interface{}) ([]datastruct.FkKV, error) {
	if fksInfo == nil {
		return nil, nil
	}

	fksValue, err := fkSliceValue(fksInfo, fks)
	if err != nil {
		return nil, errors.Wrap(err, "twf.getFKSlice")
	}

	kvs := make([]datastruct.FkKV, 0, fksValue.Len())
	for k := 0; k < fksValue.Len(); k++ {
		kvs = append(kvs, getFKKV(fksValue.Index(k), fksInfo))
	}

	return kvs, nil
}

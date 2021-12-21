package twf

import (
	"errors"
	"fmt"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"strings"
)

var (
	ErrFksIndexNotExists = errors.New("fks index not exists")
	ErrFksMustBeSLice    = errors.New("fks must me a slice")
)

func getFKValue(fksInfo *datastruct.FkInfo, originalValue interface{},
	fks ...interface{}) (fkValue interface{}, value interface{}, err error) {

	if fksInfo == nil {
		return nil, originalValue, nil
	}

	if len(fks) <= fksInfo.FksIndex {
		return "", "", ErrFksIndexNotExists
	}
	fksSlice := fks[fksInfo.FksIndex]
	if reflect.TypeOf(fksSlice).Kind() != reflect.Slice {
		return "", "", ErrFksMustBeSLice
	}
	fksValue := reflect.ValueOf(fksSlice)
	for k := 0; k < fksValue.Len(); k++ {
		v := fksValue.Index(k)
		fkKv := datastruct.FkKV{}
		for l := 0; l < v.NumField(); l++ {
			fv := v.Field(l)
			ft := v.Type().Field(l)
			if tag, ok := ft.Tag.Lookup("twf"); ok {
				tags := strings.Split(tag, ",")
				for _, e := range tags {
					if strings.HasPrefix(e, "name:") {
						if e[5:] == fksInfo.Name {
							fkKv.Name = fv
						}
						if e[5:] == fksInfo.ID {
							fkKv.ID = fv.Interface()
						}
					}
				}
			} else {
				return "", "", fmt.Errorf("twf tag wit name value must be in all fks structs")
			}
		}
		if originalValue == fkKv.ID {
			return fkKv.Name, originalValue, nil
		}
	}

	return nil, originalValue, nil
}

func getFKSlice(fksInfo *datastruct.FkInfo, fks ...interface{}) (kvs []datastruct.FkKV, err error) {
	kvs = make([]datastruct.FkKV, 0)

	if fksInfo == nil {
		return nil, nil
	}

	if len(fks) <= fksInfo.FksIndex {
		return nil, ErrFksIndexNotExists
	}
	fksSlice := fks[fksInfo.FksIndex]
	if reflect.TypeOf(fksSlice).Kind() != reflect.Slice {
		return nil, ErrFksMustBeSLice
	}
	fksValue := reflect.ValueOf(fksSlice)
	for k := 0; k < fksValue.Len(); k++ {
		v := fksValue.Index(k)
		fkKv := datastruct.FkKV{}
		for l := 0; l < v.NumField(); l++ {
			fv := v.Field(l)
			ft := v.Type().Field(l)
			if tag, ok := ft.Tag.Lookup("twf"); ok {
				tags := strings.Split(tag, ",")
				for _, e := range tags {
					if strings.HasPrefix(e, "name:") {
						if e[5:] == fksInfo.Name {
							fkKv.Name = fv
						}
						if e[5:] == fksInfo.ID {
							fkKv.ID = fv.Interface()
						}
					}
				}
			} else {
				return nil, fmt.Errorf("twf tag wit name value must be in all fks structs")
			}
		}

		kvs = append(kvs, fkKv)
	}

	return kvs, nil
}

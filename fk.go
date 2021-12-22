package twf

import (
	"fmt"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"strings"
)

// getFKValue returns fk value and original value
func getFKValue(fksInfo *datastruct.FkInfo, originalValue interface{},
	fks ...interface{}) (fkValue interface{}, value interface{}, err error) {

	if fksInfo == nil {
		return nil, originalValue, nil
	}

	if len(fks) <= fksInfo.FksIndex {
		return "", "", fmt.Errorf("twf.getFKValue: fks index %d not provided", fksInfo.FksIndex)
	}
	fksSlice := fks[fksInfo.FksIndex]
	if reflect.TypeOf(fksSlice).Kind() != reflect.Slice {
		return "", "", fmt.Errorf("twf.getFKValue: fks with index %d must be slice, not %s", fksInfo.FksIndex, reflect.TypeOf(fksSlice).Kind().String())
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

// getFKSlice return all fks of given fks info
func getFKSlice(fksInfo *datastruct.FkInfo, fks ...interface{}) (kvs []datastruct.FkKV, err error) {
	kvs = make([]datastruct.FkKV, 0)

	if fksInfo == nil {
		return nil, nil
	}

	if len(fks) <= fksInfo.FksIndex {
		return nil, fmt.Errorf("twf.getFKSlice: fks index %d not provided", fksInfo.FksIndex)
	}
	fksSlice := fks[fksInfo.FksIndex]
	if reflect.TypeOf(fksSlice).Kind() != reflect.Slice {
		return nil, fmt.Errorf("twf.getFKSlice: fks with index %d must be slice, not %s", fksInfo.FksIndex, reflect.TypeOf(fksSlice).Kind().String())
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

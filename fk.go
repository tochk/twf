package twf

import (
	"github.com/tochk/twf/datastruct"
	"reflect"
	"strings"
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
				return "", "", ErrParameterNotFound
			}
		}
		if originalValue == fkKv.ID {
			return fkKv.Name, originalValue, nil
		}
	}

	return nil, originalValue, nil
}

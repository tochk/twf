package twf

import (
	"errors"
	"github.com/tochk/twf/datastruct"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrParameterNotFound = errors.New("one or more parameters not found")
	ErrInvalidFkInfo     = errors.New("invalid fk info (fk must be {fk_array_index;id;name})")
	ErrInvalidType       = errors.New("item must be ptr to struct")
)

func GetFieldDescription(item interface{}) ([]datastruct.Field, error) {
	s := reflect.TypeOf(item)
	if s.Kind() != reflect.Ptr {
		return nil, ErrInvalidType
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
			if len(e) > 3 {
				if e[:3] == "fk:" {
					fkInfo := strings.Split(e[3:], ";")
					if len(fkInfo) != 3 {
						return nil, ErrInvalidFkInfo
					}
					fkID, err := strconv.Atoi(fkInfo[0])
					if err != nil {
						return nil, err
					}
					field.FkInfo = &datastruct.FkInfo{
						FksIndex: fkID,
						ID:       fkInfo[1],
						Name:     fkInfo[2],
					}
					field.Type = "select"
					continue
				}
			}
			if len(e) > 5 {
				if e[:5] == "name:" {
					field.Name = e[5:]
					continue
				}
				if e[:5] == "type:" {
					field.Type = e[5:]
					continue
				}
			}
			if len(e) > 6 {
				if e[:6] == "title:" {
					field.Title = e[6:]
					continue
				}
				if e[:6] == "value:" {
					field.Value = e[6:]
					continue
				}
			}
			switch e {
			case "is_not_creatable":
				field.IsNotCreatable = true
			case "is_not_disabled":
				field.IsNotDisabled = true
			case "is_not_editable":
				field.IsNotEditable = true
			case "is_not_required":
				field.IsNotRequired = true
			case "is_not_show_on_item":
				field.IsNotShowOnItem = true
			case "is_not_show_on_list":
				field.IsNotShowOnList = true
			case "process_parameters":
				field.ProcessParameters = true
			default:
				log.Print(e, tagContent, f)
				return nil, ErrParameterNotFound
			}
		}
		fields = append(fields, field)
	}
	return fields, nil
}

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
	ErrInvalidFkInfo     = errors.New("invalid fk info (fk must be {fk_slice_index;id;name})")
	ErrInvalidType       = errors.New("item must be pointer to struct")
)

func GetFieldDescription(item interface{}) ([]datastruct.Field, error) {
	s := reflect.TypeOf(item)
	if s.Kind() != reflect.Ptr {
		return nil, ErrInvalidType
	}
	//reflect.PtrTo(s) // TODO try to get pointer
	//	fmt.Println(reflect.PtrTo(reflect.TypeOf(t).Elem()).Kind().String())
	//	s := reflect.PtrTo(reflect.TypeOf(t).Elem()).Elem()
	//	for i := 0; i < s.NumField(); i++ {
	//		f := s.Field(i)
	//		fmt.Println(f.Type.Kind().String())
	//	} // TODO get it in another funcs - func for get item type from slice
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
			case strings.HasPrefix(e, "name:"):
				field.Name = e[5:]
			case strings.HasPrefix(e, "type:"):
				field.Type = e[5:]
			case strings.HasPrefix(e, "title:"):
				field.Title = e[6:]
			case strings.HasPrefix(e, "value:"):
				field.Value = e[6:]
			case e == "is_not_creatable":
				field.IsNotCreatable = true
			case e == "is_not_disabled":
				field.IsNotDisabled = true
			case e == "is_not_editable":
				field.IsNotEditable = true
			case e == "is_not_required":
				field.IsNotRequired = true
			case e == "is_not_show_on_item":
				field.IsNotShowOnItem = true
			case e == "is_not_show_on_list":
				field.IsNotShowOnList = true
			case e == "process_parameters":
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

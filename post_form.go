package twf

import (
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

func PostFormToStruct(item interface{}, r *http.Request) error {
	fields, err := GetFieldDescription(item)
	if err != nil {
		return err
	}
	postFormMap := map[string]string{}
	for k := range r.PostForm {
		postFormMap[k] = r.PostForm.Get(k)
	}
	s := reflect.ValueOf(item).Elem()
	for i := 0; i < s.NumField(); i++ {
		if _, ok := postFormMap[fields[i].Name]; !ok && fields[i].Type != "file" {
			continue
		}
		switch s.Field(i).Kind() {
		case reflect.String:
			s.Field(i).SetString(postFormMap[fields[i].Name])
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(postFormMap[fields[i].Name], 10, 64)
			if err != nil {
				return err
			}
			s.Field(i).SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(postFormMap[fields[i].Name], 10, 64)
			if err != nil {
				return err
			}
			s.Field(i).SetUint(v)
		case reflect.Float64, reflect.Float32:
			v, err := strconv.ParseFloat(postFormMap[fields[i].Name], 64)
			if err != nil {
				return err
			}
			s.Field(i).SetFloat(v)
		case reflect.Bool:
			s.Field(i).SetBool(postFormMap[fields[i].Name] == "on")
		case reflect.Slice:
			log.Println(s.Field(i).Type() == reflect.TypeOf([]byte{}))
			log.Println(fields[i].Type == "file")
			if s.Field(i).Type() == reflect.TypeOf([]byte{}) && fields[i].Type == "file" {
				file, _, err := r.FormFile(fields[i].Name)
				if err != nil {
					return err
				}
				data, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}
				s.Field(i).SetBytes(data)
			}
		}
	}
	return nil
}

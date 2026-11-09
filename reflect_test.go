package twf

import (
	"reflect"
	"testing"

	"github.com/tochk/twf/datastruct"
)

func Test_getFieldValue(t *testing.T) {
	x := 5
	p := &x
	var nilIface any

	tests := []struct {
		name  string
		field reflect.Value
		want  any
	}{
		{
			name: "simple_int",
			field: reflect.ValueOf(&struct {
				Test int
			}{2}).Elem().Field(0),
			want: 2,
		},
		{
			name: "simple_nil",
			field: reflect.ValueOf(&struct {
				Test *int
			}{nil}).Elem().Field(0),
			want: "<nil>",
		},
		{
			name: "simple_string",
			field: reflect.ValueOf(&struct {
				Test string
			}{"test"}).Elem().Field(0),
			want: "test",
		},
		{
			name: "bool",
			field: reflect.ValueOf(&struct {
				Test bool
			}{true}).Elem().Field(0),
			want: true,
		},
		{
			name: "non-nil pointer is dereferenced",
			field: reflect.ValueOf(&struct {
				Test *int
			}{p}).Elem().Field(0),
			want: 5,
		},
		{
			name: "nil string pointer",
			field: reflect.ValueOf(&struct {
				Test *string
			}{nil}).Elem().Field(0),
			want: "<nil>",
		},
		{
			// pinned: only one level of pointers is dereferenced
			name: "pointer to pointer is dereferenced once",
			field: reflect.ValueOf(&struct {
				Test **int
			}{&p}).Elem().Field(0),
			want: p,
		},
		{
			name: "nil interface field returns nil",
			field: reflect.ValueOf(&struct {
				Test any
			}{nilIface}).Elem().Field(0),
			want: nil,
		},
		{
			name: "interface field returns held value",
			field: reflect.ValueOf(&struct {
				Test any
			}{3}).Elem().Field(0),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldValue(tt.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldValue() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_getSliceElementPtrType(t *testing.T) {
	type item struct {
		Name string `twf:"name:name"`
	}

	tests := []struct {
		name    string
		item    any
		want    reflect.Type
		wantErr string
	}{
		{
			name:    "nil",
			item:    nil,
			wantErr: "expected slice, got nil",
		},
		{
			name:    "int",
			item:    1,
			wantErr: "expected slice, got int",
		},
		{
			name:    "struct",
			item:    item{},
			wantErr: "expected slice, got struct",
		},
		{
			name:    "array",
			item:    [1]item{},
			wantErr: "expected slice, got array",
		},
		{
			name:    "pointer to slice",
			item:    &[]item{},
			wantErr: "expected slice, got ptr",
		},
		{
			name:    "simpe_slice_err",
			item:    []string{},
			wantErr: "contains string, expected struct",
		},
		{
			name:    "simpe_slice_ptr_err",
			item:    []*string{},
			wantErr: "contains ptr, expected struct",
		},
		{
			name:    "simpe_slice_ptr_err_struct",
			item:    []*struct{}{},
			wantErr: "contains ptr, expected struct",
		},
		{
			name: "simpe_slice_struct",
			item: []struct{}{},
			want: reflect.TypeOf(&struct{}{}),
		},
		{
			name: "slice of named struct",
			item: []item{{"a"}},
			want: reflect.TypeFor[*item](),
		},
		{
			name: "nil slice of struct",
			item: []item(nil),
			want: reflect.TypeFor[*item](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSliceElementPtrType(tt.item)
			if checkErr(t, err, tt.wantErr) {
				if got != nil {
					t.Errorf("getSliceElementPtrType() got = %v on error, want nil", got)
				}
				return
			}
			if got != tt.want {
				t.Errorf("getSliceElementPtrType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFieldDescription(t *testing.T) {
	tests := []struct {
		name    string
		s       reflect.Type
		want    []datastruct.Field
		wantErr string
	}{
		// argument errors
		{
			name:    "simpe_not_ptr",
			s:       reflect.TypeOf(struct{}{}),
			wantErr: "expected ptr to struct, got struct",
		},
		{
			name:    "simpe_not_ptr_to_struct",
			s:       reflect.TypeFor[*string](),
			wantErr: "expected ptr to struct, got ptr to string",
		},
		{
			name:    "pointer to pointer to struct",
			s:       reflect.TypeOf(new(*struct{})),
			wantErr: "expected ptr to struct, got ptr to ptr",
		},

		// tag errors
		{
			name: "missing tag",
			s: reflect.TypeOf(&struct {
				A string
			}{}),
			wantErr: "twf tag must be present and not be empty",
		},
		{
			name: "only foreign tag",
			s: reflect.TypeOf(&struct {
				A string `json:"a"`
			}{}),
			wantErr: "twf tag must be present and not be empty",
		},
		{
			name: "empty tag",
			s: reflect.TypeOf(&struct {
				A string `twf:""`
			}{}),
			wantErr: "twf tag must be present and not be empty",
		},
		{
			name: "unknown option",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a,foo"`
			}{}),
			wantErr: `unexpected parameter in twf tag: "foo"`,
		},
		{
			name: "unknown option with value",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a,size:10"`
			}{}),
			wantErr: `unexpected parameter in twf tag: "size:10"`,
		},
		{
			name: "trailing comma",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a,"`
			}{}),
			wantErr: `unexpected parameter in twf tag: ""`,
		},
		{
			name: "space after comma",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a, required"`
			}{}),
			wantErr: `unexpected parameter in twf tag: " required"`,
		},
		{
			name: "options are case sensitive",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a,Required"`
			}{}),
			wantErr: `unexpected parameter in twf tag: "Required"`,
		},
		{
			name: "flag with value",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a,required:true"`
			}{}),
			wantErr: `unexpected parameter in twf tag: "required:true"`,
		},
		{
			name: "error on second field",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a"`
				B string
			}{}),
			wantErr: "twf tag must be present and not be empty",
		},
		{
			name: "unexported field",
			s: reflect.TypeOf(&struct {
				a string `twf:"name:a"`
			}{}),
			wantErr: "field a must be exported",
		},

		// fk errors
		{
			name: "fk with two parts",
			s: reflect.TypeOf(&struct {
				A int `twf:"name:a,fk:0;id"`
			}{}),
			wantErr: "invalid fk info {0;id}",
		},
		{
			name: "fk with four parts",
			s: reflect.TypeOf(&struct {
				A int `twf:"name:a,fk:0;id;name;x"`
			}{}),
			wantErr: "invalid fk info {0;id;name;x}",
		},
		{
			name: "fk without value",
			s: reflect.TypeOf(&struct {
				A int `twf:"name:a,fk:"`
			}{}),
			wantErr: "invalid fk info {}",
		},
		{
			name: "fk in README comma format",
			s: reflect.TypeOf(&struct {
				A int `twf:"name:a,fk:0,id,name"`
			}{}),
			wantErr: "invalid fk info {0}",
		},
		{
			name: "fk index is not int",
			s: reflect.TypeOf(&struct {
				A int `twf:"name:a,fk:a;id;name"`
			}{}),
			wantErr: "can't convert fk id to int",
		},
		{
			name: "fk index is negative",
			s: reflect.TypeOf(&struct {
				A int `twf:"name:a,fk:-1;id;name"`
			}{}),
			wantErr: "fk index must not be negative, got -1",
		},
		{
			name: "fk on non-comparable field",
			s: reflect.TypeOf(&struct {
				A []int `twf:"name:a,fk:0;id;name"`
			}{}),
			wantErr: "fk field A must be of comparable type, got []int",
		},
		{
			name: "fk on pointer to non-comparable field",
			s: reflect.TypeOf(&struct {
				A *[]int `twf:"name:a,fk:0;id;name"`
			}{}),
			wantErr: "fk field A must be of comparable type, got *[]int",
		},

		// success
		{
			name: "simpe_struct",
			s: reflect.TypeOf(&struct {
				Test string `twf:"name:test"`
			}{}),
			want: []datastruct.Field{{Name: "test", Type: "text"}},
		},
		{
			name: "complex_struct",
			s: reflect.TypeOf(&struct {
				ID      string `twf:"name:id,title:ID,no_create,no_edit"`
				Count   int    `twf:"name:count"`
				Enabled bool   `twf:"name:enabled"`
			}{}),
			want: []datastruct.Field{
				{Name: "id", Type: "text", Title: "ID", NoCreate: true, NoEdit: true},
				{Name: "count", Type: "number"},
				{Name: "enabled", Type: "checkbox"},
			},
		},
		{
			name: "type inference",
			s: reflect.TypeOf(&struct {
				I   int        `twf:"name:i"`
				I8  int8       `twf:"name:i8"`
				I16 int16      `twf:"name:i16"`
				I32 int32      `twf:"name:i32"`
				I64 int64      `twf:"name:i64"`
				U   uint       `twf:"name:u"`
				U8  uint8      `twf:"name:u8"`
				U16 uint16     `twf:"name:u16"`
				U32 uint32     `twf:"name:u32"`
				U64 uint64     `twf:"name:u64"`
				F32 float32    `twf:"name:f32"`
				F64 float64    `twf:"name:f64"`
				B   bool       `twf:"name:b"`
				S   string     `twf:"name:s"`
				P   *int       `twf:"name:p"` // pinned: pointers are not inferred as numbers
				Bs  []byte     `twf:"name:bs"`
				N   struct{}   `twf:"name:n"`
				Up  uintptr    `twf:"name:up"`
				C   complex128 `twf:"name:c"`
				Ifc any        `twf:"name:ifc"`
			}{}),
			want: []datastruct.Field{
				{Name: "i", Type: "number"},
				{Name: "i8", Type: "number"},
				{Name: "i16", Type: "number"},
				{Name: "i32", Type: "number"},
				{Name: "i64", Type: "number"},
				{Name: "u", Type: "number"},
				{Name: "u8", Type: "number"},
				{Name: "u16", Type: "number"},
				{Name: "u32", Type: "number"},
				{Name: "u64", Type: "number"},
				{Name: "f32", Type: "number"},
				{Name: "f64", Type: "number"},
				{Name: "b", Type: "checkbox"},
				{Name: "s", Type: "text"},
				{Name: "p", Type: "text"},
				{Name: "bs", Type: "text"},
				{Name: "n", Type: "text"},
				{Name: "up", Type: "text"},
				{Name: "c", Type: "text"},
				{Name: "ifc", Type: "text"},
			},
		},
		{
			name: "all options",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:n,title:T,type:email,value:v,placeholder:P,required,no_create,no_edit,not_show_on_table,process_parameters,disabled"`
			}{}),
			want: []datastruct.Field{{
				Name:              "n",
				Title:             "T",
				Type:              "email",
				Value:             "v",
				Placeholder:       "P",
				Required:          true,
				NoCreate:          true,
				NoEdit:            true,
				NotShowOnTable:    true,
				ProcessParameters: true,
				Disabled:          true,
			}},
		},
		{
			name: "option values keep colons, spaces and html",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a,title:a:b,placeholder:Enter any login,value:<a href='/u/{id}'>Edit</a>"`
			}{}),
			want: []datastruct.Field{{
				Name:        "a",
				Title:       "a:b",
				Type:        "text",
				Placeholder: "Enter any login",
				Value:       "<a href='/u/{id}'>Edit</a>",
			}},
		},
		{
			name: "empty option values",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:,title:,value:"`
			}{}),
			want: []datastruct.Field{{Type: "text"}},
		},
		{
			name: "type overrides inferred type",
			s: reflect.TypeOf(&struct {
				B bool   `twf:"name:b,type:hidden"`
				S string `twf:"name:s,type:textarea"`
			}{}),
			want: []datastruct.Field{
				{Name: "b", Type: "hidden"},
				{Name: "s", Type: "textarea"},
			},
		},
		{
			name: "fk sets fk info and select type",
			s: reflect.TypeOf(&struct {
				G int `twf:"name:g,fk:1;id;name"`
			}{}),
			want: []datastruct.Field{
				{Name: "g", Type: "select", FkInfo: &datastruct.FkInfo{FksIndex: 1, ID: "id", Name: "name"}},
			},
		},
		{
			name: "type after fk wins",
			s: reflect.TypeOf(&struct {
				G int `twf:"name:g,fk:0;id;name,type:text"`
			}{}),
			want: []datastruct.Field{
				{Name: "g", Type: "text", FkInfo: &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "name"}},
			},
		},
		{
			name: "fk after type wins",
			s: reflect.TypeOf(&struct {
				G int `twf:"name:g,type:text,fk:0;id;name"`
			}{}),
			want: []datastruct.Field{
				{Name: "g", Type: "select", FkInfo: &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "name"}},
			},
		},
		{
			name: "last duplicated option wins",
			s: reflect.TypeOf(&struct {
				A string `twf:"name:a,title:first,title:second"`
			}{}),
			want: []datastruct.Field{{Name: "a", Title: "second", Type: "text"}},
		},
		{
			name: "fk on pointer field",
			s: reflect.TypeOf(&struct {
				G *int `twf:"name:g,fk:0;id;name"`
			}{}),
			want: []datastruct.Field{
				{Name: "g", Type: "select", FkInfo: &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "name"}},
			},
		},
		{
			name: "struct without fields",
			s:    reflect.TypeOf(&struct{}{}),
			want: []datastruct.Field{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldDescription(tt.s)
			if checkErr(t, err, tt.wantErr) {
				if got != nil {
					t.Errorf("getFieldDescription() got = %v on error, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldDescription()\ngot  = %+v\nwant = %+v", got, tt.want)
			}
		})
	}
}

func Test_structPtrValue(t *testing.T) {
	type item struct {
		A string `twf:"name:a"`
	}
	var nilItem *item
	x := 1

	tests := []struct {
		name    string
		item    any
		wantErr string
	}{
		{name: "nil", item: nil, wantErr: "fn: expected ptr to struct, got nil"},
		{name: "struct value", item: item{}, wantErr: "fn: expected ptr to struct, got struct"},
		{name: "pointer to int", item: &x, wantErr: "fn: expected ptr to struct, got ptr to int"},
		{name: "typed nil pointer", item: nilItem, wantErr: "fn: expected ptr to struct, got nil ptr"},
		{name: "pointer to struct", item: &item{A: "a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := structPtrValue("fn", tt.item)
			if checkErr(t, err, tt.wantErr) {
				if got.IsValid() {
					t.Errorf("structPtrValue() returned valid value on error")
				}
				return
			}
			// the returned value must be the pointed struct itself (settable), not a copy
			if !got.CanSet() || got.Addr().Interface() != tt.item {
				t.Errorf("structPtrValue() must return the settable pointed struct")
			}
		})
	}
}

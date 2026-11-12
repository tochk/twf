package twf

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// formRequest returns request with already parsed post form
func formRequest(values url.Values) *http.Request {
	if values == nil {
		values = url.Values{}
	}
	return &http.Request{PostForm: values}
}

// multipartRequest returns unparsed multipart request with given fields and files
func multipartRequest(t *testing.T, fields map[string]string, files map[string][]byte) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	for k, v := range files {
		fw, err := w.CreateFormFile(k, k+".bin")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write(v); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodPost, "/", body)
	r.Header.Set("Content-Type", w.FormDataContentType())
	return r
}

func TestPostFormToStruct(t *testing.T) {
	type defaultTest struct {
		Name string `twf:"name:name,value:not_me"`
	}
	type kinds struct {
		S   string  `twf:"name:s"`
		I   int     `twf:"name:i"`
		I8  int8    `twf:"name:i8"`
		I16 int16   `twf:"name:i16"`
		I32 int32   `twf:"name:i32"`
		I64 int64   `twf:"name:i64"`
		U   uint    `twf:"name:u"`
		U8  uint8   `twf:"name:u8"`
		U16 uint16  `twf:"name:u16"`
		U32 uint32  `twf:"name:u32"`
		U64 uint64  `twf:"name:u64"`
		F32 float32 `twf:"name:f32"`
		F64 float64 `twf:"name:f64"`
		B   bool    `twf:"name:b"`
	}
	type text struct {
		S string `twf:"name:s"`
		I int    `twf:"name:i"`
	}
	type checkbox struct {
		B bool `twf:"name:b"`
	}
	type checkboxes struct {
		Plain    bool `twf:"name:plain"`
		Disabled bool `twf:"name:disabled,disabled"`
		NoEdit   bool `twf:"name:no_edit,no_edit"`
		NoCreate bool `twf:"name:no_create,no_create"`
		Hidden   bool `twf:"name:hidden,type:hidden"`
	}
	type file struct {
		Name   string `twf:"name:name"`
		Avatar []byte `twf:"name:avatar,type:file"`
	}
	type ignored struct {
		Raw  []byte   `twf:"name:raw"`
		List []string `twf:"name:list"`
		P    *string  `twf:"name:p"`
	}
	type stringFile struct {
		S string `twf:"name:s,type:file"`
	}
	type untagged struct {
		S string
	}

	str := "keep"

	tests := []struct {
		name string
		item any
		r    func(t *testing.T) *http.Request
		want any
	}{
		{
			name: "default values are not affecting real values",
			item: &defaultTest{},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"name": {"me"}}) },
			want: &defaultTest{Name: "me"},
		},
		{
			name: "all kinds at their bounds",
			item: &kinds{},
			r: func(*testing.T) *http.Request {
				return formRequest(url.Values{
					"s":   {"text"},
					"i":   {"-42"},
					"i8":  {"-128"},
					"i16": {"32767"},
					"i32": {"-2147483648"},
					"i64": {"9223372036854775807"},
					"u":   {"42"},
					"u8":  {"255"},
					"u16": {"65535"},
					"u32": {"4294967295"},
					"u64": {"18446744073709551615"},
					"f32": {"1.5"},
					"f64": {"-2.25"},
					"b":   {"on"},
				})
			},
			want: &kinds{
				S: "text", I: -42, I8: -128, I16: 32767, I32: -2147483648, I64: 9223372036854775807,
				U: 42, U8: 255, U16: 65535, U32: 4294967295, U64: 18446744073709551615,
				F32: 1.5, F64: -2.25, B: true,
			},
		},
		{
			name: "missing keys leave fields untouched",
			item: &text{S: "keep", I: 5},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"other": {"x"}}) },
			want: &text{S: "keep", I: 5},
		},
		{
			name: "empty value clears string",
			item: &text{S: "keep", I: 5},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"s": {""}}) },
			want: &text{S: "", I: 5},
		},
		{
			name: "first of multiple values is used",
			item: &text{},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"s": {"a", "b"}, "i": {"1", "2"}}) },
			want: &text{S: "a", I: 1},
		},
		{
			// pinned: caller must parse form before calling PostFormToStruct
			name: "unparsed form changes nothing, checkboxes included",
			item: &checkboxes{Plain: true, Disabled: true, NoEdit: true, NoCreate: true, Hidden: true},
			r: func(*testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", strings.NewReader("plain=on"))
			},
			want: &checkboxes{Plain: true, Disabled: true, NoEdit: true, NoCreate: true, Hidden: true},
		},

		// bools
		{
			name: "bool is true only for on",
			item: &checkbox{},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"b": {"on"}}) },
			want: &checkbox{B: true},
		},
		{
			name: "bool true value is not on",
			item: &checkbox{B: true},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"b": {"true"}}) },
			want: &checkbox{B: false},
		},
		{
			name: "bool 1 is not on",
			item: &checkbox{B: true},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"b": {"1"}}) },
			want: &checkbox{B: false},
		},
		{
			name: "bool empty value",
			item: &checkbox{B: true},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"b": {""}}) },
			want: &checkbox{B: false},
		},
		{
			name: "absent checkboxes",
			item: &checkboxes{Plain: true, Disabled: true, NoEdit: true, NoCreate: true, Hidden: true},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{}) },
			// unchecked checkbox is not submitted by browsers, so it is set to false.
			// disabled and no_edit fields are not submitted either and must be kept.
			// no_create fields are rendered in edit form, so they are reset as plain checkboxes.
			// non-checkbox bool is a regular input, absent means not rendered.
			want: &checkboxes{Plain: false, Disabled: true, NoEdit: true, NoCreate: false, Hidden: true},
		},
		{
			name: "checked checkboxes",
			item: &checkboxes{},
			r: func(*testing.T) *http.Request {
				return formRequest(url.Values{
					"plain": {"on"}, "disabled": {"on"}, "no_edit": {"on"}, "no_create": {"on"}, "hidden": {"on"},
				})
			},
			want: &checkboxes{Plain: true, Disabled: true, NoEdit: true, NoCreate: true, Hidden: true},
		},

		// files
		{
			name: "file is read into bytes",
			item: &file{},
			r: func(t *testing.T) *http.Request {
				r := multipartRequest(t, map[string]string{"name": "bob"}, map[string][]byte{"avatar": []byte("png")})
				if err := r.ParseMultipartForm(1 << 20); err != nil {
					t.Fatal(err)
				}
				return r
			},
			want: &file{Name: "bob", Avatar: []byte("png")},
		},
		{
			name: "missing file leaves bytes untouched",
			item: &file{Avatar: []byte("old")},
			r: func(t *testing.T) *http.Request {
				r := multipartRequest(t, map[string]string{"name": "bob"}, nil)
				if err := r.ParseMultipartForm(1 << 20); err != nil {
					t.Fatal(err)
				}
				return r
			},
			want: &file{Name: "bob", Avatar: []byte("old")},
		},
		{
			name: "file is read from unparsed multipart request",
			item: &file{Name: "keep"},
			r: func(t *testing.T) *http.Request {
				return multipartRequest(t, map[string]string{"name": "bob"}, map[string][]byte{"avatar": []byte("png")})
			},
			// form values are snapshotted before FormFile parses the body
			want: &file{Name: "keep", Avatar: []byte("png")},
		},
		{
			name: "empty file",
			item: &file{Avatar: []byte("old")},
			r: func(t *testing.T) *http.Request {
				return multipartRequest(t, nil, map[string][]byte{"avatar": {}})
			},
			want: &file{Avatar: []byte{}},
		},

		// pinned behavior
		{
			name: "unsupported types are ignored even if present",
			item: &ignored{Raw: []byte("old"), List: []string{"old"}, P: &str},
			r: func(*testing.T) *http.Request {
				return formRequest(url.Values{"raw": {"new"}, "list": {"new"}, "p": {"new"}})
			},
			want: &ignored{Raw: []byte("old"), List: []string{"old"}, P: &str},
		},
		{
			name: "string with type file is reset when absent",
			item: &stringFile{S: "keep"},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{}) },
			want: &stringFile{S: ""},
		},
		{
			name: "string with type file is set from form value",
			item: &stringFile{},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"s": {"v"}}) },
			want: &stringFile{S: "v"},
		},
		{
			name: "struct without fields",
			item: &struct{}{},
			r:    func(*testing.T) *http.Request { return formRequest(url.Values{"s": {"v"}}) },
			want: &struct{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PostFormToStruct(tt.item, tt.r(t))
			checkErr(t, err, "")
			if !reflect.DeepEqual(tt.item, tt.want) {
				t.Errorf("PostFormToStruct()\ngot  = %+v\nwant = %+v", tt.item, tt.want)
			}
		})
	}

	type number struct {
		I   int     `twf:"name:i"`
		I8  int8    `twf:"name:i8"`
		I16 int16   `twf:"name:i16"`
		I32 int32   `twf:"name:i32"`
		I64 int64   `twf:"name:i64"`
		U   uint    `twf:"name:u"`
		U8  uint8   `twf:"name:u8"`
		U16 uint16  `twf:"name:u16"`
		U32 uint32  `twf:"name:u32"`
		F32 float32 `twf:"name:f32"`
		F64 float64 `twf:"name:f64"`
	}
	var nilItem *text
	x := 1

	errTests := []struct {
		name    string
		item    any
		r       *http.Request
		wantErr string
	}{
		{name: "nil item", item: nil, r: formRequest(nil), wantErr: "twf.PostFormToStruct: expected ptr to struct, got nil"},
		{name: "struct value", item: text{}, r: formRequest(nil), wantErr: "twf.PostFormToStruct: expected ptr to struct, got struct"},
		{name: "typed nil pointer", item: nilItem, r: formRequest(nil), wantErr: "twf.PostFormToStruct: expected ptr to struct, got nil ptr"},
		{name: "pointer to int", item: &x, r: formRequest(nil), wantErr: "twf.PostFormToStruct: expected ptr to struct, got ptr to int"},
		{name: "nil request", item: &text{}, r: nil, wantErr: "twf.PostFormToStruct: request is nil"},
		{name: "bad tag", item: &untagged{}, r: formRequest(nil), wantErr: "twf.PostFormToStruct: twf tag must be present"},

		{name: "int invalid", item: &number{}, r: formRequest(url.Values{"i": {"abc"}}), wantErr: `parsing "abc": invalid syntax`},
		{name: "int empty", item: &number{}, r: formRequest(url.Values{"i": {""}}), wantErr: `parsing "": invalid syntax`},
		{name: "int float value", item: &number{}, r: formRequest(url.Values{"i": {"1.5"}}), wantErr: `parsing "1.5": invalid syntax`},
		{name: "int8 overflow", item: &number{}, r: formRequest(url.Values{"i8": {"128"}}), wantErr: `parsing "128": value out of range`},
		{name: "int8 underflow", item: &number{}, r: formRequest(url.Values{"i8": {"-129"}}), wantErr: `parsing "-129": value out of range`},
		{name: "int16 overflow", item: &number{}, r: formRequest(url.Values{"i16": {"32768"}}), wantErr: "value out of range"},
		{name: "int32 overflow", item: &number{}, r: formRequest(url.Values{"i32": {"2147483648"}}), wantErr: "value out of range"},
		{name: "int64 overflow", item: &number{}, r: formRequest(url.Values{"i64": {"9223372036854775808"}}), wantErr: "value out of range"},
		{name: "uint negative", item: &number{}, r: formRequest(url.Values{"u": {"-1"}}), wantErr: `parsing "-1": invalid syntax`},
		{name: "uint8 overflow", item: &number{}, r: formRequest(url.Values{"u8": {"256"}}), wantErr: `parsing "256": value out of range`},
		{name: "uint16 overflow", item: &number{}, r: formRequest(url.Values{"u16": {"65536"}}), wantErr: "value out of range"},
		{name: "uint32 overflow", item: &number{}, r: formRequest(url.Values{"u32": {"4294967296"}}), wantErr: "value out of range"},
		{name: "float invalid", item: &number{}, r: formRequest(url.Values{"f64": {"x"}}), wantErr: `parsing "x": invalid syntax`},
		{name: "float32 overflow", item: &number{}, r: formRequest(url.Values{"f32": {"1e39"}}), wantErr: `parsing "1e39": value out of range`},
		{name: "float64 overflow", item: &number{}, r: formRequest(url.Values{"f64": {"1e309"}}), wantErr: "value out of range"},
		{
			name: "file from non-multipart request",
			item: &file{},
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("name=bob"))
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r
			}(),
			wantErr: "twf.PostFormToStruct: request Content-Type isn't multipart/form-data",
		},
	}
	for _, tt := range errTests {
		t.Run("error: "+tt.name, func(t *testing.T) {
			checkErr(t, PostFormToStruct(tt.item, tt.r), tt.wantErr)
		})
	}
}

func TestPostFormToStructEmptySelect(t *testing.T) {
	type item struct {
		FkID int    `twf:"name:fk_id,fk:0;id;name"`
		Code string `twf:"name:code,type:select"`
		Num  int    `twf:"name:num"`
	}
	tests := []struct {
		name    string
		values  url.Values
		want    item
		wantErr string
	}{
		{name: "empty select resets to zero", values: url.Values{"fk_id": {""}, "code": {""}}, want: item{Num: 7}},
		{name: "selected value is parsed", values: url.Values{"fk_id": {"3"}, "code": {"x"}}, want: item{FkID: 3, Code: "x", Num: 7}},
		{name: "empty non-select number still fails", values: url.Values{"num": {""}}, wantErr: `parsing "": invalid syntax`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := item{FkID: 5, Code: "old", Num: 7}
			err := PostFormToStruct(&got, formRequest(tt.values))
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("PostFormToStruct() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("PostFormToStruct() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostFormToStruct() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

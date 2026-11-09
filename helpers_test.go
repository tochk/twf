package twf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/tochk/twf/datastruct"
)

// recorder captures everything the page builders pass to TWF render functions
type recorder struct {
	listFields []datastruct.Field
	items      [][]interface{}
	formFields []datastruct.Field
	selectKVs  [][]datastruct.FkKV
	selected   []interface{}
}

// newRecorder returns TWF with deterministic marker render functions, so tests can assert
// exact output of page builders without depending on real HTML templates
func newRecorder() (*TWF, *recorder) {
	rec := &recorder{}
	formItem := func(prefix string) func(datastruct.Field) string {
		return func(f datastruct.Field) string {
			rec.formFields = append(rec.formFields, f)
			return prefix + "(" + f.Name + "=" + f.Value + ")"
		}
	}
	t := &TWF{
		HeadFunc:   func(title string) string { return "H(" + title + ")" },
		MenuFunc:   func() string { return "M" },
		FooterFunc: func() string { return "F" },
		ListFunc: func(fields []datastruct.Field, content string) string {
			rec.listFields = fields
			names := make([]string, 0, len(fields))
			for _, f := range fields {
				names = append(names, f.Name)
			}
			return "L[" + strings.Join(names, ",") + "]{" + content + "}"
		},
		ListItemFunc: func(values []interface{}) string {
			rec.items = append(rec.items, values)
			parts := make([]string, 0, len(values))
			for _, v := range values {
				parts = append(parts, fmt.Sprint(v))
			}
			return "R[" + strings.Join(parts, "|") + "]"
		},
		FormFunc:         func(link, content string) string { return "FORM(" + link + "){" + content + "}" },
		FormItemText:     formItem("T"),
		FormItemTextarea: formItem("A"),
		FormItemCheckbox: formItem("C"),
		FormItemSelect: func(f datastruct.Field, kvs []datastruct.FkKV, selected interface{}) string {
			rec.formFields = append(rec.formFields, f)
			rec.selectKVs = append(rec.selectKVs, kvs)
			rec.selected = append(rec.selected, selected)
			return fmt.Sprintf("S(%s=%s;%s;%v)", f.Name, f.Value, kvString(kvs), selected)
		},
	}
	return t, rec
}

// kvString renders kvs as [id:name,...]. FkKV.Name holds reflect.Value, fmt prints its underlying value
func kvString(kvs []datastruct.FkKV) string {
	if kvs == nil {
		return "nil"
	}
	parts := make([]string, 0, len(kvs))
	for _, kv := range kvs {
		parts = append(parts, fmt.Sprintf("%v:%v", kv.ID, kv.Name))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// unwrapValues replaces reflect.Value elements (fk names) with their underlying values
func unwrapValues(values []interface{}) []interface{} {
	if values == nil {
		return nil
	}
	res := make([]interface{}, 0, len(values))
	for _, v := range values {
		res = append(res, unwrapValue(v))
	}
	return res
}

func unwrapValue(v interface{}) interface{} {
	if rv, ok := v.(reflect.Value); ok {
		return rv.Interface()
	}
	return v
}

// checkErr fails the test if err doesn't match wantErr. Empty wantErr means no error is expected,
// otherwise the error message must contain wantErr. Returns true if an error was expected
func checkErr(t *testing.T, err error, wantErr string) bool {
	t.Helper()
	if wantErr == "" {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return false
	}
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", wantErr)
	}
	if !strings.Contains(err.Error(), wantErr) {
		t.Fatalf("expected error containing %q, got %q", wantErr, err.Error())
	}
	return true
}

func intPtr(i int) *int { return &i }

// fk test fixtures shared between test files
type fkGroup struct {
	ID   int    `twf:"name:id,title:ID"`
	Name string `twf:"title:Name,name:name"`
}

type fkOther struct {
	Code string `twf:"name:code"`
}

type fkUntagged struct {
	ID   int `twf:"name:id"`
	Name string
}

type fkUnexported struct {
	id   int    `twf:"name:id"`
	Name string `twf:"name:name"`
}

var fkGroupInfo = &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "name"}

var fkGroups = []fkGroup{{1, "admins"}, {2, "users"}}

// formErrorCase describes invalid input shared by AddForm and EditForm tests
type formErrorCase struct {
	name    string
	item    interface{}
	fks     []interface{}
	wantErr string
}

func formErrorCases(fn string) []formErrorCase {
	type item struct {
		A string `twf:"name:a"`
	}
	type untagged struct {
		A string
	}
	type withFk struct {
		G int `twf:"name:g,fk:0;id;name"`
	}
	type hiddenFk struct {
		G int `twf:"name:g,fk:1;id;name,no_create,no_edit"`
	}
	var nilItem *item
	x := 1

	return []formErrorCase{
		{name: "nil item", item: nil, wantErr: fn + ": expected ptr to struct, got nil"},
		{name: "struct value", item: item{}, wantErr: fn + ": expected ptr to struct, got struct"},
		{name: "typed nil pointer", item: nilItem, wantErr: fn + ": expected ptr to struct, got nil ptr"},
		{name: "pointer to int", item: &x, wantErr: fn + ": expected ptr to struct, got ptr to int"},
		{name: "bad tag", item: &untagged{}, wantErr: fn + ": twf tag must be present and not be empty"},
		{
			name: "unexported field",
			item: &struct {
				a string `twf:"name:a"`
			}{},
			wantErr: fn + ": twf.getFieldDescription: field a must be exported",
		},
		{name: "fk index not provided", item: &withFk{}, wantErr: fn + ": twf.getFKSlice: fks index 0 not provided"},
		{
			// pinned: fk misconfiguration is reported even for fields not rendered in the form
			name:    "fk index not provided for hidden field",
			item:    &hiddenFk{},
			fks:     []interface{}{fkGroups},
			wantErr: fn + ": twf.getFKSlice: fks index 1 not provided",
		},
		{
			name:    "nil fks element",
			item:    &item{},
			fks:     []interface{}{nil},
			wantErr: fn + ": twf.validateFk: fks with index 0 must be slice, got nil",
		},
		{
			name:    "unused invalid fks element",
			item:    &item{},
			fks:     []interface{}{fkGroups, 5},
			wantErr: fn + ": twf.validateFk: fks with index 1 must be slice, not int",
		},
		{
			name:    "fks slice of int",
			item:    &withFk{},
			fks:     []interface{}{[]int{1}},
			wantErr: "fks with index 0 must be slice of structs, got slice of int",
		},
	}
}

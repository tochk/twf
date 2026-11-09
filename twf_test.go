package twf

import (
	"reflect"
	"testing"

	"github.com/tochk/twf/twftemplates"
)

func TestNew(t *testing.T) {
	tw := New()

	v := reflect.ValueOf(tw).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			t.Errorf("New().%s is nil", v.Type().Field(i).Name)
		}
	}

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"HeadFunc", tw.HeadFunc, twftemplates.Head},
		{"MenuFunc", tw.MenuFunc, twftemplates.Menu},
		{"ListFunc", tw.ListFunc, twftemplates.ItemList},
		{"ListItemFunc", tw.ListItemFunc, twftemplates.Item},
		{"FooterFunc", tw.FooterFunc, twftemplates.Footer},
		{"FormFunc", tw.FormFunc, twftemplates.Form},
		{"FormItemText", tw.FormItemText, twftemplates.FormItem},
		{"FormItemTextarea", tw.FormItemTextarea, twftemplates.FormItemTextarea},
		{"FormItemSelect", tw.FormItemSelect, twftemplates.FormItemSelect},
		{"FormItemCheckbox", tw.FormItemCheckbox, twftemplates.FormItemCheckbox},
	}
	if len(tests) != v.NumField() {
		t.Fatalf("TWF has %d fields, test checks %d: add new fields to this test", v.NumField(), len(tests))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if reflect.ValueOf(tt.got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("New().%s is not wired to the expected default template", tt.name)
			}
		})
	}
}

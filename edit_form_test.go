package twf

import (
	"reflect"
	"testing"

	"github.com/tochk/twf/datastruct"
)

func TestTWF_EditForm(t *testing.T) {
	type dispatch struct {
		S string `twf:"name:s"`
		N int    `twf:"name:n"`
		P string `twf:"name:p,type:password"`
		F []byte `twf:"name:f,type:file"`
		B bool   `twf:"name:b"`
		T string `twf:"name:t,type:textarea"`
	}
	type selects struct {
		G int    `twf:"name:g,fk:0;id;name"`
		S string `twf:"name:s,type:select"`
	}
	type ptrFk struct {
		G *int `twf:"name:g,fk:0;id;name"`
	}
	type values struct {
		Def  string `twf:"name:def,value:default"`
		Zero string `twf:"name:zero,value:default"`
		N    int    `twf:"name:n,value:5"`
		B    bool   `twf:"name:b,value:true"`
		P    *int   `twf:"name:p"`
		PV   *int   `twf:"name:pv"`
	}
	type flags struct {
		ID     int    `twf:"name:id,no_edit"`
		Create string `twf:"name:create,no_create"`
		Hidden string `twf:"name:hidden,not_show_on_table"`
	}
	type pp struct {
		ID    int    `twf:"name:id,no_edit"`
		Login string `twf:"name:login"`
		Link  string `twf:"name:link,value:/u/{id}/{login},process_parameters"`
		Chain string `twf:"name:chain,value:[{link}],process_parameters"`
		Later string `twf:"name:later,value:{after},process_parameters"`
		Empty string `twf:"name:empty,process_parameters"`
		After string `twf:"name:after"`
	}
	type ppRefs struct {
		P   *int   `twf:"name:p,no_edit"`
		G   int    `twf:"name:g,fk:0;id;name"`
		Ref string `twf:"name:ref,value:{p}|{g},process_parameters"`
	}

	tests := []struct {
		name         string
		item         interface{}
		fks          []interface{}
		want         string
		wantSelected []interface{}
	}{
		{
			name: "dispatch by type",
			item: &dispatch{S: "x", N: 7, P: "secret", F: []byte("ab"), B: true, T: "long"},
			want: "H(Edit)MFORM(/edit){T(s=x)T(n=7)T(p=secret)T(f=[97 98])C(b=true)A(t=long)}F",
		},
		{
			name:         "select gets raw field value as selected",
			item:         &selects{G: 2, S: "x"},
			fks:          []interface{}{fkGroups},
			want:         "H(Edit)MFORM(/edit){S(g=2;[1:admins,2:users];2)S(s=x;nil;x)}F",
			wantSelected: []interface{}{2, "x"},
		},
		{
			name:         "select on pointer field gets dereferenced value",
			item:         &ptrFk{G: intPtr(1)},
			fks:          []interface{}{fkGroups},
			want:         "H(Edit)MFORM(/edit){S(g=1;[1:admins,2:users];1)}F",
			wantSelected: []interface{}{1},
		},
		{
			// pinned: nil pointer is rendered and selected as "<nil>" string
			name:         "select on nil pointer field",
			item:         &ptrFk{},
			fks:          []interface{}{fkGroups},
			want:         "H(Edit)MFORM(/edit){S(g=<nil>;[1:admins,2:users];<nil>)}F",
			wantSelected: []interface{}{"<nil>"},
		},
		{
			name: "value tag never overrides real values, including zero values",
			item: &values{Def: "real", PV: intPtr(5)},
			want: "H(Edit)MFORM(/edit){T(def=real)T(zero=)T(n=0)C(b=false)T(p=<nil>)T(pv=5)}F",
		},
		{
			name: "no_edit is skipped, other flags are rendered",
			item: &flags{ID: 1, Create: "c", Hidden: "h"},
			want: "H(Edit)MFORM(/edit){T(create=c)T(hidden=h)}F",
		},
		{
			name: "process parameters use earlier fields including no_edit",
			item: &pp{ID: 7, Login: "bob", Link: "ignored", Chain: "ignored", Empty: "e", After: "a"},
			want: "H(Edit)MFORM(/edit){T(login=bob)T(link=/u/7/bob)T(chain=[/u/7/bob])T(later={after})T(empty=e)T(after=a)}F",
		},
		{
			name: "process parameters use <nil> pointer and raw fk id",
			item: &ppRefs{G: 1},
			fks:  []interface{}{fkGroups},
			want: "H(Edit)MFORM(/edit){S(g=1;[1:admins,2:users];1)T(ref=<nil>|1)}F",
		},
		{
			name: "struct without fields",
			item: &struct{}{},
			want: "H(Edit)MFORM(/edit){}F",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw, rec := newRecorder()
			got, err := tw.EditForm("Edit", tt.item, "/edit", tt.fks...)
			checkErr(t, err, "")
			if got != tt.want {
				t.Errorf("EditForm()\ngot  = %s\nwant = %s", got, tt.want)
			}
			if tt.wantSelected != nil && !reflect.DeepEqual(rec.selected, tt.wantSelected) {
				t.Errorf("FormItemSelect selected = %#v, want %#v", rec.selected, tt.wantSelected)
			}
		})
	}

	for _, tt := range formErrorCases("twf.EditForm") {
		t.Run("error: "+tt.name, func(t *testing.T) {
			tw, _ := newRecorder()
			got, err := tw.EditForm("Edit", tt.item, "/edit", tt.fks...)
			checkErr(t, err, tt.wantErr)
			if got != "" {
				t.Errorf("EditForm() = %q on error, want empty", got)
			}
		})
	}
}

func TestTWF_EditForm_fieldMetadata(t *testing.T) {
	type meta struct {
		Email string `twf:"name:email,title:E-mail,type:email,placeholder:Enter email,required,disabled"`
		G     int    `twf:"name:g,title:Group,fk:0;id;name"`
	}
	tw, rec := newRecorder()
	_, err := tw.EditForm("Edit", &meta{Email: "a@b", G: 1}, "/edit", fkGroups)
	checkErr(t, err, "")

	want := []datastruct.Field{
		{Name: "email", Title: "E-mail", Type: "email", Placeholder: "Enter email", Required: true, Disabled: true, Value: "a@b"},
		{Name: "g", Title: "Group", Type: "select", Value: "1", FkInfo: &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "name"}},
	}
	if !reflect.DeepEqual(rec.formFields, want) {
		t.Errorf("form item fields\ngot  = %+v\nwant = %+v", rec.formFields, want)
	}
	if kvs := kvString(rec.selectKVs[0]); kvs != "[1:admins,2:users]" {
		t.Errorf("FormItemSelect kvs = %s, want [1:admins,2:users]", kvs)
	}
}

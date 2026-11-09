package twf

import (
	"reflect"
	"testing"

	"github.com/tochk/twf/datastruct"
)

func TestTWF_AddForm(t *testing.T) {
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
	type values struct {
		Def  string `twf:"name:def,value:default"`
		Real string `twf:"name:real"`
		P    *int   `twf:"name:p"`
		PV   *int   `twf:"name:pv"`
		B    bool   `twf:"name:b,value:true"`
		F    bool   `twf:"name:f"`
	}
	type flags struct {
		ID     int    `twf:"name:id,no_create"`
		Edit   string `twf:"name:edit,no_edit"`
		Hidden string `twf:"name:hidden,not_show_on_table"`
	}
	type pp struct {
		ID    int    `twf:"name:id,no_create"`
		Login string `twf:"name:login"`
		Link  string `twf:"name:link,value:/u/{id}/{login},process_parameters"`
		Chain string `twf:"name:chain,value:[{link}],process_parameters"`
		Later string `twf:"name:later,value:{after},process_parameters"`
		Empty string `twf:"name:empty,process_parameters"`
		After string `twf:"name:after"`
	}
	type ppRefs struct {
		P   *int   `twf:"name:p,no_create"`
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
			item: &dispatch{S: "x", N: 7, P: "secret", B: true, T: "long"},
			want: "H(Add)MFORM(/add){T(s=x)T(n=7)T(p=secret)T(f=[])C(b=true)A(t=long)}F",
		},
		{
			name:         "select never has selected value",
			item:         &selects{G: 2, S: "x"},
			fks:          []interface{}{fkGroups},
			want:         "H(Add)MFORM(/add){S(g=2;[1:admins,2:users];<nil>)S(s=x;nil;<nil>)}F",
			wantSelected: []interface{}{nil, nil},
		},
		{
			name: "value tag is a default, nil pointer is empty",
			item: &values{Def: "real value", Real: "r", PV: intPtr(5)},
			want: "H(Add)MFORM(/add){T(def=default)T(real=r)T(p=)T(pv=5)C(b=true)C(f=false)}F",
		},
		{
			name: "no_create is skipped, other flags are rendered",
			item: &flags{ID: 1, Edit: "e", Hidden: "h"},
			want: "H(Add)MFORM(/add){T(edit=e)T(hidden=h)}F",
		},
		{
			name: "process parameters use earlier fields including no_create",
			item: &pp{ID: 7, Login: "bob", Link: "ignored", Chain: "ignored", Empty: "e", After: "a"},
			want: "H(Add)MFORM(/add){T(login=bob)T(link=/u/7/bob)T(chain=[/u/7/bob])T(later={after})T(empty=e)T(after=a)}F",
		},
		{
			name: "process parameters use empty nil pointer and raw fk id",
			item: &ppRefs{G: 1},
			fks:  []interface{}{fkGroups},
			want: "H(Add)MFORM(/add){S(g=1;[1:admins,2:users];<nil>)T(ref=|1)}F",
		},
		{
			name: "struct without fields",
			item: &struct{}{},
			want: "H(Add)MFORM(/add){}F",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw, rec := newRecorder()
			got, err := tw.AddForm("Add", tt.item, "/add", tt.fks...)
			checkErr(t, err, "")
			if got != tt.want {
				t.Errorf("AddForm()\ngot  = %s\nwant = %s", got, tt.want)
			}
			if tt.wantSelected != nil && !reflect.DeepEqual(rec.selected, tt.wantSelected) {
				t.Errorf("FormItemSelect selected = %#v, want %#v", rec.selected, tt.wantSelected)
			}
		})
	}

	for _, tt := range formErrorCases("twf.AddForm") {
		t.Run("error: "+tt.name, func(t *testing.T) {
			tw, _ := newRecorder()
			got, err := tw.AddForm("Add", tt.item, "/add", tt.fks...)
			checkErr(t, err, tt.wantErr)
			if got != "" {
				t.Errorf("AddForm() = %q on error, want empty", got)
			}
		})
	}
}

func TestTWF_AddForm_fieldMetadata(t *testing.T) {
	type meta struct {
		Email string `twf:"name:email,title:E-mail,type:email,placeholder:Enter email,required,disabled"`
		G     int    `twf:"name:g,title:Group,fk:0;id;name"`
	}
	tw, rec := newRecorder()
	_, err := tw.AddForm("Add", &meta{Email: "a@b", G: 1}, "/add", fkGroups)
	checkErr(t, err, "")

	want := []datastruct.Field{
		{Name: "email", Title: "E-mail", Type: "email", Placeholder: "Enter email", Required: true, Disabled: true, Value: "a@b"},
		{Name: "g", Title: "Group", Type: "select", Value: "1", FkInfo: &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "name"}},
	}
	if !reflect.DeepEqual(rec.formFields, want) {
		t.Errorf("form item fields\ngot  = %+v\nwant = %+v", rec.formFields, want)
	}
}

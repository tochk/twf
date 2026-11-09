package twf

import (
	"reflect"
	"testing"

	"github.com/tochk/twf/datastruct"
)

func Test_processParameters(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		fields map[string]string
		want   string
	}{
		{name: "no placeholders", value: "static", fields: map[string]string{"id": "1"}, want: "static"},
		{name: "one placeholder", value: "/users/{id}", fields: map[string]string{"id": "1"}, want: "/users/1"},
		{name: "repeated placeholder", value: "{id}-{id}", fields: map[string]string{"id": "7"}, want: "7-7"},
		{
			name:   "several placeholders",
			value:  `<a href="/u/{id}/{login}">{login}</a>`,
			fields: map[string]string{"id": "1", "login": "bob"},
			want:   `<a href="/u/1/bob">bob</a>`,
		},
		{name: "unknown placeholder kept", value: "/u/{uid}", fields: map[string]string{"id": "1"}, want: "/u/{uid}"},
		{name: "empty value", value: "", fields: map[string]string{"id": "1"}, want: ""},
		{name: "nil fields", value: "/u/{id}", fields: nil, want: "/u/{id}"},
		{name: "empty braces kept", value: "{}", fields: map[string]string{"id": "1"}, want: "{}"},
		{name: "placeholder without braces is not replaced", value: "id", fields: map[string]string{"id": "1"}, want: "id"},
		{name: "replacement with empty string", value: "a{id}b", fields: map[string]string{"id": ""}, want: "ab"},
		{
			// substituted values must not be expanded again, otherwise result depends on map order
			name:   "single pass",
			value:  "{a}-{b}",
			fields: map[string]string{"a": "{b}", "b": "x"},
			want:   "{b}-x",
		},
		{
			name:   "single pass reverse",
			value:  "{a}-{b}",
			fields: map[string]string{"a": "x", "b": "{a}"},
			want:   "x-{a}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// repeat to catch map iteration order dependency
			for i := 0; i < 20; i++ {
				if got := processParameters(tt.value, tt.fields); got != tt.want {
					t.Fatalf("processParameters() = %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func Test_filterTableFields(t *testing.T) {
	a := datastruct.Field{Name: "a"}
	b := datastruct.Field{Name: "b", NotShowOnTable: true}
	c := datastruct.Field{Name: "c"}

	tests := []struct {
		name   string
		fields []datastruct.Field
		want   []datastruct.Field
	}{
		{name: "nil", fields: nil, want: []datastruct.Field{}},
		{name: "none hidden", fields: []datastruct.Field{a, c}, want: []datastruct.Field{a, c}},
		{name: "all hidden", fields: []datastruct.Field{b, b}, want: []datastruct.Field{}},
		{name: "mixed keeps order", fields: []datastruct.Field{c, b, a}, want: []datastruct.Field{c, a}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterTableFields(tt.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterTableFields() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_generateTableItemsSlice(t *testing.T) {
	type defaultTest struct {
		Name string `twf:"name:name,value:not_me"`
	}
	type defaultTest2 struct {
		ID   int    `twf:"name:id,not_show_on_table"`
		Name string `twf:"name:name,value:me{id},process_parameters"`
	}
	type defaultTest3 struct {
		Name string `twf:"name:name,fk:0;name;name"`
	}
	type defaultFks3 struct {
		Name string `twf:"name:name"`
	}
	type plain struct {
		ID     int    `twf:"name:id"`
		Name   string `twf:"name:name"`
		Active bool   `twf:"name:active"`
	}
	type pointers struct {
		P *int    `twf:"name:p"`
		N *string `twf:"name:n"`
	}
	type ppMulti struct {
		ID    int    `twf:"name:id"`
		Login string `twf:"name:login"`
		Link  string `twf:"name:link,value:/u/{id}/{login},process_parameters"`
	}
	type ppLater struct {
		Link string `twf:"name:link,value:x{id},process_parameters"`
		ID   int    `twf:"name:id"`
	}
	type ppNil struct {
		P    *int   `twf:"name:p"`
		Link string `twf:"name:link,value:{p},process_parameters"`
	}
	type ppEmpty struct {
		S string `twf:"name:s,process_parameters"`
	}
	type withFk struct {
		G int `twf:"name:g,fk:0;id;name"`
	}
	type withFk64 struct {
		G int64 `twf:"name:g,fk:0;id;name"`
	}
	type withFkPtr struct {
		G *int `twf:"name:g,fk:0;id;name"`
	}
	type ppFk struct {
		G    int    `twf:"name:g,fk:0;id;name"`
		Link string `twf:"name:link,value:/g/{g},process_parameters"`
	}
	type ppHiddenFk struct {
		G    int    `twf:"name:g,fk:0;id;name,not_show_on_table"`
		Link string `twf:"name:link,value:/g/{g},process_parameters"`
	}

	tests := []struct {
		name    string
		slice   interface{}
		i       int
		fks     []interface{}
		want    []interface{}
		wantErr string
	}{
		{
			name:  "default values should not affect real values",
			slice: []defaultTest{{Name: "me"}},
			want:  []interface{}{"me"},
		},
		{
			name:  "process parameters should work properly",
			slice: []defaultTest2{{ID: 1, Name: "me"}},
			want:  []interface{}{"me1"},
		},
		{
			name:  "fks",
			slice: []defaultTest3{{Name: "name"}},
			fks:   []interface{}{[]defaultFks3{{"test"}, {"name"}}},
			want:  []interface{}{"name"},
		},
		{
			name:  "plain values keep their types",
			slice: []plain{{1, "a", true}},
			want:  []interface{}{1, "a", true},
		},
		{
			name:  "pointers",
			slice: []pointers{{P: intPtr(5)}},
			want:  []interface{}{5, "<nil>"},
		},
		{
			name:  "process parameters with several fields",
			slice: []ppMulti{{ID: 7, Login: "bob", Link: "ignored"}},
			want:  []interface{}{7, "bob", "/u/7/bob"},
		},
		{
			// pinned: only fields declared before the process_parameters field are available
			name:  "process parameters can't use later fields",
			slice: []ppLater{{ID: 1}},
			want:  []interface{}{"x{id}", 1},
		},
		{
			name:  "process parameters with nil pointer",
			slice: []ppNil{{}},
			want:  []interface{}{"<nil>", "<nil>"},
		},
		{
			name:  "process parameters without value is empty",
			slice: []ppEmpty{{S: "real"}},
			want:  []interface{}{""},
		},
		{
			name:  "fk resolved to name",
			slice: []withFk{{G: 2}},
			fks:   []interface{}{fkGroups},
			want:  []interface{}{"users"},
		},
		{
			name:  "fk not found shows raw value",
			slice: []withFk{{G: 99}},
			fks:   []interface{}{fkGroups},
			want:  []interface{}{99},
		},
		{
			name:  "fk type mismatch shows raw value",
			slice: []withFk64{{G: 1}},
			fks:   []interface{}{fkGroups},
			want:  []interface{}{int64(1)},
		},
		{
			name:  "fk on pointer field",
			slice: []withFkPtr{{G: intPtr(1)}},
			fks:   []interface{}{fkGroups},
			want:  []interface{}{"admins"},
		},
		{
			name:  "fk on nil pointer field",
			slice: []withFkPtr{{}},
			fks:   []interface{}{fkGroups},
			want:  []interface{}{"<nil>"},
		},
		{
			// pinned: process parameters receive fk name, not id
			name:  "process parameters use resolved fk name",
			slice: []ppFk{{G: 1}},
			fks:   []interface{}{fkGroups},
			want:  []interface{}{"admins", "/g/admins"},
		},
		{
			name:  "hidden fk still feeds process parameters",
			slice: []ppHiddenFk{{G: 2}},
			fks:   []interface{}{fkGroups},
			want:  []interface{}{"/g/users"},
		},
		{
			name:  "row index selects row",
			slice: []plain{{1, "a", true}, {2, "b", false}},
			i:     1,
			want:  []interface{}{2, "b", false},
		},
		{
			name:    "fk error propagates",
			slice:   []withFk{{G: 1}},
			fks:     nil,
			wantErr: "fks index 0 not provided",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := getSliceElementPtrType(tt.slice)
			checkErr(t, err, "")
			fields, err := getFieldDescription(item)
			checkErr(t, err, "")

			got, err := generateTableItemsSlice(reflect.ValueOf(tt.slice), fields, tt.i, tt.fks)
			if checkErr(t, err, tt.wantErr) {
				return
			}
			if got := unwrapValues(got); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateTableItemsSlice() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestTWF_Table(t *testing.T) {
	type defaultTest3 struct {
		Name string `twf:"name:name,fk:0;name;name"`
	}
	type defaultFks3 struct {
		Name string `twf:"name:name"`
	}
	type row struct {
		ID   int    `twf:"name:id"`
		Name string `twf:"name:name"`
	}
	type hidden struct {
		ID   int    `twf:"name:id,title:ID,not_show_on_table"`
		Name string `twf:"name:name,title:Name"`
		Link string `twf:"name:link,value:/u/{id},process_parameters"`
	}
	type withFk struct {
		ID int `twf:"name:id"`
		G  int `twf:"name:g,fk:0;id;name"`
	}
	type withFk1 struct {
		G int `twf:"name:g,fk:1;id;name"`
	}
	type untagged struct {
		Name string
	}

	tests := []struct {
		name      string
		slice     interface{}
		fks       []interface{}
		want      string
		wantList  []string
		wantItems int
		wantErr   string
	}{
		{
			name:      "fks no panic",
			slice:     []defaultTest3{{"test"}},
			fks:       []interface{}{[]defaultFks3{{"test"}, {"name"}}},
			want:      "H(Users)ML[name]{R[test]}F",
			wantList:  []string{"name"},
			wantItems: 1,
		},
		{
			name:      "empty slice renders header only",
			slice:     []row{},
			want:      "H(Users)ML[id,name]{}F",
			wantList:  []string{"id", "name"},
			wantItems: 0,
		},
		{
			name:      "nil slice of struct renders header only",
			slice:     []row(nil),
			want:      "H(Users)ML[id,name]{}F",
			wantList:  []string{"id", "name"},
			wantItems: 0,
		},
		{
			name:      "rows in order",
			slice:     []row{{1, "a"}, {2, "b"}},
			want:      "H(Users)ML[id,name]{R[1|a]R[2|b]}F",
			wantList:  []string{"id", "name"},
			wantItems: 2,
		},
		{
			name:      "hidden fields are not listed but feed process parameters",
			slice:     []hidden{{ID: 3, Name: "c"}},
			want:      "H(Users)ML[name,link]{R[c|/u/3]}F",
			wantList:  []string{"name", "link"},
			wantItems: 1,
		},
		{
			name:      "fk resolved and unused fks allowed",
			slice:     []withFk{{1, 2}, {2, 5}},
			fks:       []interface{}{fkGroups, []fkOther{}},
			want:      "H(Users)ML[id,g]{R[1|users]R[2|5]}F",
			wantList:  []string{"id", "g"},
			wantItems: 2,
		},
		{
			// pinned: fk indexes are only checked while rendering rows
			name:      "fk misconfiguration with empty slice",
			slice:     []withFk1{},
			want:      "H(Users)ML[g]{}F",
			wantList:  []string{"g"},
			wantItems: 0,
		},

		{name: "nil slice", slice: nil, wantErr: "twf.Table: twf.getSliceElementPtrType: expected slice, got nil"},
		{name: "struct", slice: row{}, wantErr: "twf.Table: twf.getSliceElementPtrType: expected slice, got struct"},
		{name: "pointer to slice", slice: &[]row{}, wantErr: "expected slice, got ptr"},
		{name: "slice of int", slice: []int{1}, wantErr: "slice contains int, expected struct"},
		{name: "slice of pointers", slice: []*row{{1, "a"}}, wantErr: "slice contains ptr, expected struct"},
		{name: "bad tag", slice: []untagged{{"a"}}, wantErr: "twf.Table: twf tag must be present"},
		{
			name:    "nil fks element",
			slice:   []row{},
			fks:     []interface{}{nil},
			wantErr: "twf.Table: twf.validateFk: fks with index 0 must be slice, got nil",
		},
		{
			name:    "non-slice fks element",
			slice:   []row{},
			fks:     []interface{}{fkGroups, 5},
			wantErr: "fks with index 1 must be slice, not int",
		},
		{
			name:    "fks slice of int",
			slice:   []row{},
			fks:     []interface{}{[]int{1}},
			wantErr: "must be slice of structs, got slice of int",
		},
		{
			name:    "fks struct without tags",
			slice:   []row{},
			fks:     []interface{}{[]fkUntagged{}},
			wantErr: "must be present in all fks structs",
		},
		{
			name:    "fks struct with unexported field",
			slice:   []row{},
			fks:     []interface{}{[]fkUnexported{}},
			wantErr: "field id must be exported",
		},
		{
			name:    "fk index out of range",
			slice:   []withFk1{{1}},
			fks:     []interface{}{fkGroups},
			wantErr: "twf.Table: twf.getFKValue: fks index 1 not provided",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw, rec := newRecorder()
			got, err := tw.Table("Users", tt.slice, tt.fks...)
			if checkErr(t, err, tt.wantErr) {
				if got != "" {
					t.Errorf("Table() = %q on error, want empty", got)
				}
				return
			}
			if got != tt.want {
				t.Errorf("Table()\ngot  = %s\nwant = %s", got, tt.want)
			}
			gotList := make([]string, 0, len(rec.listFields))
			for _, f := range rec.listFields {
				gotList = append(gotList, f.Name)
			}
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("ListFunc fields = %v, want %v", gotList, tt.wantList)
			}
			if len(rec.items) != tt.wantItems {
				t.Errorf("ListItemFunc called %d times, want %d", len(rec.items), tt.wantItems)
			}
		})
	}
}

func TestTWF_Table_listFieldsMetadata(t *testing.T) {
	type row struct {
		ID   int    `twf:"name:id,title:ID,not_show_on_table"`
		Name string `twf:"name:name,title:Name,placeholder:p,required"`
	}
	tw, rec := newRecorder()
	_, err := tw.Table("Users", []row{{1, "a"}})
	checkErr(t, err, "")

	want := []datastruct.Field{{Name: "name", Title: "Name", Type: "text", Placeholder: "p", Required: true}}
	if !reflect.DeepEqual(rec.listFields, want) {
		t.Errorf("ListFunc fields = %+v, want %+v", rec.listFields, want)
	}
}

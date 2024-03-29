package twf

import (
	"fmt"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"testing"
)

func Test_generateTableItemsSlice(t *testing.T) {
	type defaultTest struct {
		Name string `twf:"name:name,value:not_me"`
	}

	slice := []defaultTest{{Name: "me"}}

	item, err := getSliceElementPtrType(slice)
	if err != nil {
		t.Fatalf("getSliceElementPtrType() error = %v", err)
	}

	fields, err := getFieldDescription(item)
	if err != nil {
		t.Fatalf("getFieldDescription() error = %v", err)
	}

	s := reflect.ValueOf(slice)

	type defaultTest2 struct {
		ID   int    `twf:"name:id,not_show_on_table"`
		Name string `twf:"name:name,value:me{id},process_parameters"`
	}

	slice2 := []defaultTest2{{ID: 1, Name: "me"}}

	item2, err := getSliceElementPtrType(slice2)
	if err != nil {
		t.Fatalf("getSliceElementPtrType() error = %v", err)
	}

	fields2, err := getFieldDescription(item2)
	if err != nil {
		t.Fatalf("getFieldDescription() error = %v", err)
	}

	s2 := reflect.ValueOf(slice2)

	type defaultTest3 struct {
		Name string `twf:"name:name,fk:0;name;name"`
	}

	type defaultFks3 struct {
		Name string `twf:"name:name"`
	}

	slice3 := []defaultTest3{{Name: "name"}}

	item3, err := getSliceElementPtrType(slice3)
	if err != nil {
		t.Fatalf("getSliceElementPtrType() error = %v", err)
	}

	fields3, err := getFieldDescription(item3)
	if err != nil {
		t.Fatalf("getFieldDescription() error = %v", err)
	}

	s3 := reflect.ValueOf(slice3)

	type args struct {
		s      reflect.Value
		fields []datastruct.Field
		i      int
		fks    []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "default values should not affect real values",
			args: args{
				s:      s,
				fields: fields,
				i:      0,
				fks:    nil,
			},
			want:    []interface{}{"me"},
			wantErr: false,
		},
		{
			name: "process parameters should work properly",
			args: args{
				s:      s2,
				fields: fields2,
				i:      0,
				fks:    nil,
			},
			want:    []interface{}{"me1"},
			wantErr: false,
		},
		{
			name: "fks",
			args: args{
				s:      s3,
				fields: fields3,
				i:      0,
				fks:    []interface{}{[]defaultFks3{{"test"}, {"name"}}},
			},
			want:    []interface{}{"name"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateTableItemsSlice(tt.args.s, tt.args.fields, tt.args.i, tt.args.fks)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateTableItemsSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(fmt.Sprintf("%+v", got), fmt.Sprintf("%+v", tt.want)) {
				t.Errorf("generateTableItemsSlice() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestTWF_Table(t1 *testing.T) {
	type args struct {
		title string
		slice interface{}
		fks   []interface{}
	}

	type defaultTest3 struct {
		Name string `twf:"name:name,fk:0;name;name"`
	}

	type defaultFks3 struct {
		Name string `twf:"name:name"`
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fks no panic",
			args: args{
				title: "Test",
				slice: []defaultTest3{{"test"}},
				fks:   []interface{}{[]defaultFks3{{"test"}, {"name"}}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := New()
			_, err := t.Table(tt.args.title, tt.args.slice, tt.args.fks...)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Table() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

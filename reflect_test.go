package twf

import (
	"github.com/tochk/twf/datastruct"
	"reflect"
	"testing"
)

func Test_getFieldValue(t *testing.T) {
	type args struct {
		field reflect.Value
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "simple_int",
			args: args{field: reflect.ValueOf(&struct {
				Test int
			}{2}).Elem().Field(0)},
			want: 2,
		},
		{
			name: "simple_nil",
			args: args{field: reflect.ValueOf(&struct {
				Test *int
			}{nil}).Elem().Field(0)},
			want: "<nil>",
		},
		{
			name: "simple_string",
			args: args{field: reflect.ValueOf(&struct {
				Test string
			}{"test"}).Elem().Field(0)},
			want: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldValue(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSliceElementPtrType(t *testing.T) {
	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Type
		wantErr bool
	}{
		{
			name:    "simpe_slice_err",
			args:    args{[]string{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "simpe_slice_ptr_err",
			args:    args{[]*string{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "simpe_slice_ptr_err_struct",
			args:    args{[]*struct{}{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "simpe_slice_struct",
			args:    args{[]struct{}{}},
			want:    reflect.TypeOf(&struct{}{}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSliceElementPtrType(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSliceElementPtrType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSliceElementPtrType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFieldDescription(t *testing.T) {
	type args struct {
		s reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    []datastruct.Field
		wantErr bool
	}{
		{
			name:    "simpe_not_ptr",
			args:    args{reflect.TypeOf(struct{}{})},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "simpe_not_ptr_to_struct",
			args:    args{reflect.TypeOf(new(string))},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simpe_struct",
			args: args{reflect.TypeOf(&struct {
				Test string `twf:"name:test"`
			}{
				Test: "",
			})},
			want:    []datastruct.Field{{Name: "test", Type: "text"}},
			wantErr: false,
		},
		{
			name: "complex_struct",
			args: args{reflect.TypeOf(&struct {
				ID      string `twf:"name:id,title:ID,no_create,no_edit"`
				Count   int    `twf:"name:count"`
				Enabled bool   `twf:"name:enabled"`
			}{
				ID: "",
			})},
			want: []datastruct.Field{
				{Name: "id", Type: "text", Title: "ID", NoCreate: true, NoEdit: true},
				{Name: "count", Type: "number"},
				{Name: "enabled", Type: "checkbox"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldDescription(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFieldDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldDescription() got = %v, want %v", got, tt.want)
			}
		})
	}
}

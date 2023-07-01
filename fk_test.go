package twf

import (
	"fmt"
	"github.com/tochk/twf/datastruct"
	"reflect"
	"testing"
)

func Test_getFKValue(t *testing.T) {
	type defaultFks struct {
		Name string `twf:"name:name"`
	}
	type args struct {
		fksInfo       *datastruct.FkInfo
		originalValue interface{}
		fks           []interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantFkValue interface{}
		wantValue   interface{}
		wantErr     bool
	}{
		{
			name: "default",
			args: args{
				fksInfo: &datastruct.FkInfo{
					FksIndex: 0,
					ID:       "name",
					Name:     "name",
				},
				originalValue: interface{}("test"),
				fks: []interface{}{
					[]defaultFks{{"test"}, {"test2"}},
				},
			},
			wantFkValue: "test",
			wantValue:   "test",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFkValue, gotValue, err := getFKValue(tt.args.fksInfo, tt.args.originalValue, tt.args.fks)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFKValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(fmt.Sprintf("%+v", gotFkValue), fmt.Sprintf("%+v", tt.wantFkValue)) {
				t.Errorf("getFKValue() gotFkValue = %v, want %v", gotFkValue, tt.wantFkValue)
			}
			if !reflect.DeepEqual(fmt.Sprintf("%+v", gotValue), fmt.Sprintf("%+v", tt.wantValue)) {
				t.Errorf("getFKValue() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

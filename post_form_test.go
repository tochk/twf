package twf

import (
	"net/http"
	"testing"
)

func TestPostFormToStruct(t *testing.T) {
	type defaultTest struct {
		Name string `twf:"name:name,value:not_me"`
	}

	type args struct {
		item interface{}
		r    *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantName string
	}{
		{
			name: "default values are not affecting real values",
			args: args{
				item: &defaultTest{},
				r: &http.Request{
					PostForm: map[string][]string{
						"name": {"me"},
					},
				},
			},
			wantName: "me",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PostFormToStruct(tt.args.item, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PostFormToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.item.(*defaultTest).Name != tt.wantName {
				t.Fatalf("PostFormToStruct() wantName = %v, got %v", tt.args.item.(*defaultTest).Name, tt.wantName)
			}
		})
	}
}

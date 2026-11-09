package twf

import (
	"reflect"
	"testing"

	"github.com/tochk/twf/datastruct"
)

// fkErrorCase describes invalid fks input, shared by getFKValue and getFKSlice tests
type fkErrorCase struct {
	name    string
	info    *datastruct.FkInfo
	fks     []interface{}
	wantErr string
}

func fkErrorCases() []fkErrorCase {
	return []fkErrorCase{
		{
			name:    "no fks",
			info:    fkGroupInfo,
			fks:     nil,
			wantErr: "fks index 0 not provided",
		},
		{
			name:    "index out of range",
			info:    &datastruct.FkInfo{FksIndex: 1, ID: "id", Name: "name"},
			fks:     []interface{}{fkGroups},
			wantErr: "fks index 1 not provided",
		},
		{
			name:    "negative index",
			info:    &datastruct.FkInfo{FksIndex: -1, ID: "id", Name: "name"},
			fks:     []interface{}{fkGroups},
			wantErr: "fks index -1 not provided",
		},
		{
			name:    "nil fks element",
			info:    fkGroupInfo,
			fks:     []interface{}{nil},
			wantErr: "fks with index 0 must be slice, got nil",
		},
		{
			name:    "fks element is not slice",
			info:    fkGroupInfo,
			fks:     []interface{}{fkGroup{1, "admins"}},
			wantErr: "fks with index 0 must be slice, not struct",
		},
		{
			name:    "fks element is slice of int",
			info:    fkGroupInfo,
			fks:     []interface{}{[]int{1}},
			wantErr: "fks with index 0 must be slice of structs, got slice of int",
		},
		{
			name:    "fks element is slice of pointers",
			info:    fkGroupInfo,
			fks:     []interface{}{[]*fkGroup{{1, "admins"}}},
			wantErr: "fks with index 0 must be slice of structs, got slice of ptr",
		},
		{
			name:    "fks struct field without twf tag",
			info:    fkGroupInfo,
			fks:     []interface{}{[]fkUntagged{{1, "admins"}}},
			wantErr: "must be present in all fks structs, missing on field Name",
		},
		{
			name:    "empty fks slice of struct without twf tag",
			info:    fkGroupInfo,
			fks:     []interface{}{[]fkUntagged{}},
			wantErr: "must be present in all fks structs, missing on field Name",
		},
		{
			name:    "fks struct with unexported field",
			info:    fkGroupInfo,
			fks:     []interface{}{[]fkUnexported{{1, "admins"}}},
			wantErr: "fks with index 0: field id must be exported",
		},
		{
			name:    "selected fks element is invalid",
			info:    &datastruct.FkInfo{FksIndex: 1, ID: "id", Name: "name"},
			fks:     []interface{}{fkGroups, 5},
			wantErr: "fks with index 1 must be slice, not int",
		},
	}
}

func Test_getFKValue(t *testing.T) {
	type defaultFks struct {
		Name string `twf:"name:name"`
	}

	tests := []struct {
		name          string
		info          *datastruct.FkInfo
		originalValue interface{}
		fks           []interface{}
		wantFkValue   interface{}
		wantValue     interface{}
	}{
		{
			name: "default",
			info: &datastruct.FkInfo{
				FksIndex: 0,
				ID:       "name",
				Name:     "name",
			},
			originalValue: "test",
			fks: []interface{}{
				[]defaultFks{{"test"}, {"test2"}},
			},
			wantFkValue: "test",
			wantValue:   "test",
		},
		{
			name:          "nil fk info returns original value",
			info:          nil,
			originalValue: 1,
			fks:           nil,
			wantFkValue:   nil,
			wantValue:     1,
		},
		{
			name:          "found by id, name tag is not first option",
			info:          fkGroupInfo,
			originalValue: 2,
			fks:           []interface{}{fkGroups},
			wantFkValue:   "users",
			wantValue:     2,
		},
		{
			name:          "first match wins on duplicated ids",
			info:          fkGroupInfo,
			originalValue: 1,
			fks:           []interface{}{[]fkGroup{{1, "first"}, {1, "second"}}},
			wantFkValue:   "first",
			wantValue:     1,
		},
		{
			name:          "selects fks slice by index",
			info:          &datastruct.FkInfo{FksIndex: 1, ID: "id", Name: "name"},
			originalValue: 1,
			fks:           []interface{}{[]fkOther{{"x"}}, fkGroups},
			wantFkValue:   "admins",
			wantValue:     1,
		},
		{
			name:          "not found",
			info:          fkGroupInfo,
			originalValue: 3,
			fks:           []interface{}{fkGroups},
			wantFkValue:   nil,
			wantValue:     3,
		},
		{
			name:          "empty fks slice",
			info:          fkGroupInfo,
			originalValue: 1,
			fks:           []interface{}{[]fkGroup{}},
			wantFkValue:   nil,
			wantValue:     1,
		},
		{
			name:          "id type mismatch is not a match",
			info:          fkGroupInfo,
			originalValue: int64(1),
			fks:           []interface{}{fkGroups},
			wantFkValue:   nil,
			wantValue:     int64(1),
		},
		{
			name:          "nil pointer value is not a match",
			info:          fkGroupInfo,
			originalValue: "<nil>",
			fks:           []interface{}{fkGroups},
			wantFkValue:   nil,
			wantValue:     "<nil>",
		},
		{
			name:          "id field missing in fks struct",
			info:          &datastruct.FkInfo{FksIndex: 0, ID: "uid", Name: "name"},
			originalValue: 1,
			fks:           []interface{}{fkGroups},
			wantFkValue:   nil,
			wantValue:     1,
		},
		{
			// pinned: matched element without name field is reported as not found
			name:          "name field missing in fks struct",
			info:          &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "title"},
			originalValue: 1,
			fks:           []interface{}{fkGroups},
			wantFkValue:   nil,
			wantValue:     1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFkValue, gotValue, err := getFKValue(tt.info, tt.originalValue, tt.fks)
			checkErr(t, err, "")
			if got := unwrapValue(gotFkValue); !reflect.DeepEqual(got, tt.wantFkValue) {
				t.Errorf("getFKValue() gotFkValue = %#v, want %#v", got, tt.wantFkValue)
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("getFKValue() gotValue = %#v, want %#v", gotValue, tt.wantValue)
			}
		})
	}

	for _, tt := range fkErrorCases() {
		t.Run("error: "+tt.name, func(t *testing.T) {
			_, _, err := getFKValue(tt.info, 1, tt.fks)
			checkErr(t, err, "twf.getFKValue")
			checkErr(t, err, tt.wantErr)
		})
	}
}

func Test_getFKSlice(t *testing.T) {
	t.Run("nil fk info", func(t *testing.T) {
		kvs, err := getFKSlice(nil, fkGroups)
		checkErr(t, err, "")
		if kvs != nil {
			t.Errorf("getFKSlice() = %v, want nil", kvs)
		}
	})

	tests := []struct {
		name   string
		info   *datastruct.FkInfo
		fks    []interface{}
		want   string
		wantID []interface{}
	}{
		{
			name:   "all elements in order",
			info:   fkGroupInfo,
			fks:    []interface{}{fkGroups},
			want:   "[1:admins,2:users]",
			wantID: []interface{}{1, 2},
		},
		{
			name:   "selects fks slice by index",
			info:   &datastruct.FkInfo{FksIndex: 1, ID: "code", Name: "code"},
			fks:    []interface{}{fkGroups, []fkOther{{"a"}, {"b"}}},
			want:   "[a:a,b:b]",
			wantID: []interface{}{"a", "b"},
		},
		{
			name:   "empty fks slice",
			info:   fkGroupInfo,
			fks:    []interface{}{[]fkGroup{}},
			want:   "[]",
			wantID: []interface{}{},
		},
		{
			name:   "name field missing",
			info:   &datastruct.FkInfo{FksIndex: 0, ID: "id", Name: "title"},
			fks:    []interface{}{fkGroups[:1]},
			want:   "[1:<nil>]",
			wantID: []interface{}{1},
		},
		{
			name:   "id field missing",
			info:   &datastruct.FkInfo{FksIndex: 0, ID: "uid", Name: "name"},
			fks:    []interface{}{fkGroups[:1]},
			want:   "[<nil>:admins]",
			wantID: []interface{}{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kvs, err := getFKSlice(tt.info, tt.fks...)
			checkErr(t, err, "")
			if got := kvString(kvs); got != tt.want {
				t.Errorf("getFKSlice() = %s, want %s", got, tt.want)
			}
			gotID := make([]interface{}, 0, len(kvs))
			for _, kv := range kvs {
				gotID = append(gotID, kv.ID)
			}
			// ids keep their original types, they are compared with field values
			if !reflect.DeepEqual(gotID, tt.wantID) {
				t.Errorf("getFKSlice() ids = %#v, want %#v", gotID, tt.wantID)
			}
		})
	}

	for _, tt := range fkErrorCases() {
		t.Run("error: "+tt.name, func(t *testing.T) {
			kvs, err := getFKSlice(tt.info, tt.fks...)
			checkErr(t, err, "twf.getFKSlice")
			checkErr(t, err, tt.wantErr)
			if kvs != nil {
				t.Errorf("getFKSlice() = %v on error, want nil", kvs)
			}
		})
	}
}

func Test_validateFks(t *testing.T) {
	tests := []struct {
		name    string
		fks     []interface{}
		wantErr string
	}{
		{name: "nil", fks: nil},
		{name: "valid fks", fks: []interface{}{fkGroups, []fkOther{}}},
		{
			name:    "reports index of first invalid element",
			fks:     []interface{}{fkGroups, []fkOther{}, []int{}, nil},
			wantErr: "fks with index 2 must be slice of structs, got slice of int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkErr(t, validateFks(tt.fks), tt.wantErr)
		})
	}
}

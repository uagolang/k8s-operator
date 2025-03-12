package validator

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetErrors(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wantErrs Errors
	}{
		{
			name: "no validation errors",
			args: args{
				err: errors.New("internal error"),
			},
			wantErrs: Errors{},
		},
		{
			name: "custom validation errors",
			args: args{
				err: Errors{
					{Field: "foo", Message: "bar"},
				},
			},
			wantErrs: Errors{
				{Field: "foo", Message: "bar"},
			},
		},
		{
			name: "custom validation error",
			args: args{
				err: Error{Field: "foo", Message: "bar"},
			},
			wantErrs: Errors{
				{Field: "foo", Message: "bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := GetErrors(tt.args.err); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("GetErrors() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

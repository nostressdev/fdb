package utils

import (
	"testing"

	"github.com/nostressdev/fdb/errors"
)

func TestValidate(t *testing.T) {
	type args struct {
		expression bool
		text       string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "expression is true",
			args: args{
				expression: true,
				text:       "test",
			},
			wantPanic: true,
		},
		{
			name: "expression is false",
			args: args{
				expression: false,
				text:       "test",
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.wantPanic && errors.GetType(r.(error)) == errors.ValidationError {
						return
					}
					t.Fatalf("Validate() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			Validate(tt.args.expression, tt.args.text)
		})
	}
}

func TestValidatef(t *testing.T) {
	type args struct {
		expression bool
		format     string
		args       []interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "expression is true",
			args: args{
				expression: true,
				format:     "%s %s",
				args:       []interface{}{"test", "test"},
			},
			wantPanic: true,
		},
		{
			name: "expression is false",
			args: args{
				expression: false,
				format:     "%s %s",
				args:       []interface{}{"test", "test"},
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.wantPanic && errors.GetType(r.(error)) == errors.ValidationError {
						return
					}
					t.Fatalf("Validate() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			Validatef(tt.args.expression, tt.args.format, tt.args.args...)
		})
	}
}

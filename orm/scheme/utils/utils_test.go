package utils

import (
	"github.com/nostressdev/fdb/utils/errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
					if errors.GetType(r.(error)) != errors.ValidationError {
						t.Fatalf("Validate() panic is not validation error")
					}
				}
			}()
			if !tt.wantPanic {
				assert.NotPanics(t, func() { Validate(tt.args.expression, tt.args.text) }, "Validate() should not panic")
			} else {
				Validate(tt.args.expression, tt.args.text)
			}
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
					if errors.GetType(r.(error)) != errors.ValidationError {
						t.Fatalf("Validate() panic is not validation error")
					}
				}
			}()
			if !tt.wantPanic {
				assert.NotPanics(t, func() { Validatef(tt.args.expression, tt.args.format, tt.args.args...) }, "Validate() should not panic")
			} else {
				Validatef(tt.args.expression, tt.args.format, tt.args.args...)
			}
		})
	}
}

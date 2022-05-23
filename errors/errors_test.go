package errors

import (
	"testing"

	"github.com/pkg/errors"
)

func TestErrorWrapping(t *testing.T) {
	err := errors.New("error")
	tests := []struct {
		name string
		args ErrorType
	}{
		{
			name: "parsing error",
			args: ParsingError,
		},
		{
			name: "validation error",
			args: ValidationError,
		},
		{
			name: "internal error",
			args: InternalError,
		},
		{
			name: "no type error",
			args: NoType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newErr := tt.args.Wrap(err, "wrapped error")
			oldErr := Unwrap(newErr)
			if oldErr.Error() != "wrapped error: "+err.Error() {
				t.Fatalf("error unwrapping error: got %q, want %q", oldErr.Error(), err.Error())
			}
		})
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		name string
		args error
		want ErrorType
	}{
		{
			name: "parsing error",
			args: ParsingError.New("error"),
			want: ParsingError,
		},
		{
			name: "validation error",
			args: ValidationError.New("error"),
			want: ValidationError,
		},
		{
			name: "internal error",
			args: InternalError.New("error"),
			want: InternalError,
		},
		{
			name: "explicit no type error",
			args: NoType.New("error"),
			want: NoType,
		},
		{
			name: "implicit no type error",
			args: errors.New("error"),
			want: NoType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetType(tt.args); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		args    error
		wantErr string
	}{
		{
			name:    "parsing error",
			args:    ParsingError.New("error"),
			wantErr: "parsing error: error",
		},
		{
			name:    "validation error",
			args:    ValidationError.New("error"),
			wantErr: "validation error: error",
		},
		{
			name:    "internal error",
			args:    InternalError.New("error"),
			wantErr: "internal error: error (please report this error)",
		},
		{
			name:    "no type error",
			args:    NoType.New("error"),
			wantErr: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.Error() != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", tt.args.Error(), tt.wantErr)
			}
		})
	}
}

package errors

import (
	"testing"

	"github.com/pkg/errors"
)

func TestErrorWrapping(t *testing.T) {
	err := errors.New("error")
	newErr := ParsingError.Wrap(err, "wrapped error")
	if newErr.Error() != "parsing error: wrapped error: error" {
		t.Fatalf("error wrapping error: got %q, want %q", newErr.Error(), "wrapped error: error")
	}
	oldErr := Unwrap(err)
	if oldErr != err {
		t.Fatalf("error unwrapping error: got %q, want %q", oldErr.Error(), err.Error())
	}
}

func TestAddedMessage(t *testing.T) {
	err := ParsingError.New("error")
	if err.Error() != "parsing error: error" {
		t.Fatalf("error adding message: got %q, want %q", err.Error(), "parsing error: error")
	}
	err = ValidationError.New("error")
	if err.Error() != "validation error: error" {
		t.Fatalf("error adding message: got %q, want %q", err.Error(), "validation error: error")
	}
	err = InternalError.New("error")
	if err.Error() != "internal error: error (please report this error)" {
		t.Fatalf("error adding message: got %q, want %q", err.Error(), "internal error: error")
	}
	err = NoType.New("error")
	if err.Error() != "error" {
		t.Fatalf("error adding message: got %q, want %q", err.Error(), "error")
	}
}

func TestErrorType(t *testing.T) {
	err := ParsingError.New("error")
	if GetType(err) != ParsingError {
		t.Fatalf("error type: got %q, want %q", GetType(err), ParsingError)
	}
	err = ValidationError.New("error")
	if GetType(err) != ValidationError {
		t.Fatalf("error type: got %q, want %q", GetType(err), ValidationError)
	}
	err = InternalError.New("error")
	if GetType(err) != InternalError {
		t.Fatalf("error type: got %q, want %q", GetType(err), InternalError)
	}
	err = NoType.New("error")
	if GetType(err) != NoType {
		t.Fatalf("error type: got %q, want %q", GetType(err), NoType)
	}
	err = errors.New("error")
	if GetType(err) != NoType {
		t.Fatalf("error type: got %q, want %q", GetType(err), NoType)
	}
}


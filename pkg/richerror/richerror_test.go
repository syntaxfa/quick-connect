package richerror_test

import (
	"testing"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type errFields struct {
	operation string
	wrapError error
	message   string
	kind      richerror.Kind
	meta      map[string]interface{}
}

func TestRichError(t *testing.T) {
	// define test cases
	richErr1 := richerror.New("test_one")
	richErr2 := richerror.New("test_two").WithKind(richerror.KindInvalid).WithMessage("test_message_two").
		WithWrapError(richerror.New("test_operation").WithMessage("test_message").
			WithKind(richerror.KindUnexpected))
	richErr3 := richerror.New("test_three").WithWrapError(richerror.New("test_error")).
		WithMessage("in_message").
		WithKind(richerror.KindForbidden)

	tests := []struct {
		input    richerror.RichError
		expected errFields
	}{
		{
			input: richErr1,
			expected: errFields{
				operation: "test_one",
				wrapError: nil,
				message:   "",
				kind:      0,
				meta:      nil,
			},
		},
		{
			input: richErr2,
			expected: errFields{
				operation: "test_two",
				wrapError: nil,
				message:   "test_message_two",
				kind:      richerror.KindInvalid,
				meta:      nil,
			},
		},
		{
			input: richErr3,
			expected: errFields{
				operation: "test_three",
				wrapError: nil,
				message:   "in_message",
				kind:      richerror.KindForbidden,
				meta:      nil,
			},
		},
	}

	for _, test := range tests {
		expect := test.expected
		input := test.input

		if expect.kind != input.Kind() {
			t.Fatalf("expected kind: %d, but got: %d", expect.kind, input.Kind())
		}
		if expect.operation != input.Operation() {
			t.Fatalf("expected operation: %s, but got: %s", expect.operation, input.Operation())
		}
		if expect.message != input.Message() {
			t.Fatalf("expected message: %s, but got: %s", expect.message, input.Message())
		}
	}
}

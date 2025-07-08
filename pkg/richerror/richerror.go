package richerror

import (
	"errors"
)

type Kind int

const (
	KindInvalid      Kind = iota + 1 // entity unprocessable
	KindForbidden                    // authorization error
	KindNotFound                     // record not found
	KindUnexpected                   // internal server error
	KindUnAuthorized                 // authentication error
	KindBadRequest                   // bad request
	KindConflict                     // conflict
)

type RichError struct {
	operation   string
	wrapError   error
	message     string
	kind        Kind
	meta        map[string]interface{}
	errorFields map[string]string
}

func New(op string) RichError {
	return RichError{operation: op}
}

func (r RichError) Error() string {
	if r.message == "" && r.wrapError != nil {
		return r.wrapError.Error()
	}

	return r.message
}

func (r RichError) WithMessage(message string) RichError {
	r.message = message

	return r
}

func (r RichError) WithWrapError(err error) RichError {
	r.wrapError = err

	return r
}

func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind

	return r
}

func (r RichError) WithMeta(meta map[string]interface{}) RichError {
	r.meta = meta

	return r
}

func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	var err RichError
	if errors.As(r.wrapError, &err) {
		return err.Kind()
	}

	return r.kind
}

func (r RichError) Message() string {
	return r.Error()
}

func (r RichError) Meta() map[string]interface{} {
	if r.meta != nil {
		return r.meta
	}

	var err RichError
	if errors.As(r.wrapError, &err) {
		return err.Meta()
	}

	return r.meta
}

func (r RichError) Operation() string {
	return r.operation
}

func (r RichError) WithErrorFields(fields map[string]string) RichError {
	r.errorFields = fields

	return r
}

func (r RichError) ErrorFields() map[string]string {
	return r.errorFields
}

func (r RichError) ExtraDetail() map[string]interface{} {
	details := make(map[string]interface{})
	details["op"] = r.operation
	details["message"] = r.Error()
	details["kind"] = r.Kind()
	details["meta"] = r.Meta()
	details["errorFields"] = r.ErrorFields()

	if r.wrapError != nil {
		var richErr RichError
		if errors.As(r.wrapError, &richErr) {
			details["error"] = richErr.ExtraDetail()
		} else {
			details["error"] = r.wrapError.Error()
		}
	}

	return details
}

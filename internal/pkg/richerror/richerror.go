package richerror

import (
	"fmt"
	"runtime"
)

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

func (k Kind) String() string {
	switch k {
	case KindInvalid:
		return "Invalid"
	case KindForbidden:
		return "Forbidden"
	case KindNotFound:
		return "NotFound"
	case KindUnexpected:
		return "Unexpected"
	default:
		return "Unknown"
	}
}

type Op string

type RichError struct {
	Operation Op
	Err       error
	Message   string
	Kind      Kind
	Meta      map[string]interface{}
	file      string
	line      int
}

// New creates a new RichError with the given operation
func New(op Op) *RichError {
	_, file, line, _ := runtime.Caller(1)
	return &RichError{
		Operation: op,
		file:      file,
		line:      line,
	}
}

// Wrap wraps an existing error with operation context
func Wrap(err error, op Op) *RichError {
	if err == nil {
		return nil
	}
	_, file, line, _ := runtime.Caller(1)
	return &RichError{
		Operation: op,
		Err:       err,
		file:      file,
		line:      line,
	}
}

// Builder methods - now return *RichError for efficiency
func (r *RichError) WithOp(op Op) *RichError {
	r.Operation = op
	return r
}

func (r *RichError) WithMessage(message string) *RichError {
	r.Message = message
	return r
}

func (r *RichError) WithKind(kind Kind) *RichError {
	r.Kind = kind
	return r
}

func (r *RichError) WithMeta(key string, value any) *RichError {
	if r.Meta == nil {
		r.Meta = make(map[string]any)
	}
	r.Meta[key] = value
	return r
}

func (r *RichError) WithMetaMap(meta map[string]any) *RichError {
	r.Meta = meta
	return r
}

func (r *RichError) WithErr(err error) *RichError {
	r.Err = err
	return r
}

// Error implements the error interface
func (r *RichError) Error() string {
	if r.Message != "" {
		if r.Err != nil {
			return fmt.Sprintf("%s: %v", r.Message, r.Err)
		}
		return r.Message
	}
	if r.Err != nil {
		return r.Err.Error()
	}
	return fmt.Sprintf("%s: %s", r.Operation, r.Kind)
}

// Unwrap returns the wrapped error (for errors.Is and errors.As)
func (r *RichError) Unwrap() error {
	return r.Err
}

// GetKind returns the kind, unwrapping if necessary
func (r *RichError) GetKind() Kind {
	if r.Kind != 0 {
		return r.Kind
	}
	if re, ok := r.Err.(*RichError); ok {
		return re.GetKind()
	}
	return 0
}

// GetMessage returns the message, unwrapping if necessary
func (r *RichError) GetMessage() string {
	if r.Message != "" {
		return r.Message
	}
	if re, ok := r.Err.(*RichError); ok {
		return re.GetMessage()
	}
	if r.Err != nil {
		return r.Err.Error()
	}
	return ""
}

// Location returns where the error was created
func (r *RichError) Location() string {
	return fmt.Sprintf("%s:%d", r.file, r.line)
}
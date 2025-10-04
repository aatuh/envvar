package getters

import "strings"

// ErrKind describes the class of error.
type ErrKind int

const (
	// ErrMissing is the error kind for missing values.
	ErrMissing ErrKind = iota + 1
	ErrType
)

// KeyError is an error for envvar key-related errors.
type KeyError struct {
	Key  string
	Kind ErrKind
	Msg  string
}

// Error returns the error message.
//
// Returns:
//   - string: The error message.
func (e *KeyError) Error() string {
	var b strings.Builder
	b.WriteString("envvar: ")
	switch e.Kind {
	case ErrMissing:
		b.WriteString("missing ")
	case ErrType:
		b.WriteString("type error for ")
	}
	b.WriteString(e.Key)
	if e.Msg != "" {
		b.WriteString(": ")
		b.WriteString(e.Msg)
	}
	return b.String()
}

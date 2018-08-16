package errtree

import (
	"errors"
	"strings"
)

type ErrorTree struct {
	err         error
	descendants []error
}

func NewString(msg string, descendants ...error) *ErrorTree {
	return New(errors.New(msg), descendants...)
}

func New(err error, descendants ...error) *ErrorTree {
	return &ErrorTree{
		err:         err,
		descendants: descendants,
	}
}

func (e *ErrorTree) error(sb *strings.Builder, depth int) {
	const indent = "    "
	indentation := strings.Repeat(indent, depth)
	sb.WriteString(indentation)
	sb.WriteString(e.err.Error())

	if len(e.descendants) == 0 {
		return
	}

	depth++
	sb.WriteRune(':')
	for _, err := range e.descendants {
		if err == nil {
			continue
		}
		sb.WriteString("\n")
		if et, ok := err.(*ErrorTree); ok {
			et.error(sb, depth)
			continue
		}
		sb.WriteString(indentation)
		sb.WriteString(indent)
		sb.WriteString(err.Error())
	}
}

func (e *ErrorTree) Len() (len int) {
	for _, d := range e.descendants {
		if et, ok := d.(*ErrorTree); ok {
			len += et.Len()
		} else {
			len += 1
		}
	}
	return len + 1
}

func (e *ErrorTree) Add(errs ...error) {
	e.descendants = append(e.descendants, errs...)
}

func (e *ErrorTree) Error() string {
	var sb strings.Builder
	e.error(&sb, 0)
	return sb.String()
}

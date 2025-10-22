package ws

import (
	"fmt"
	"strings"
)

// ConnError represents an error related to a specific connection.
// It includes the connection ID, the operation being performed, and the underlying error.
type ConnError struct {
	ConnID ID
	Op     string
	Err    error
}

func (e *ConnError) Error() string {
	return fmt.Sprintf("conn[%s] %s: %v", e.ConnID, e.Op, e.Err)
}

func (e *ConnError) Unwrap() error {
	return e.Err
}

type ConnErrors []ConnError

func (errs ConnErrors) Error() string {
	if len(errs) == 0 {
		return "no connection errors"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d connection error(s):\n", len(errs)))
	for _, e := range errs {
		sb.WriteString(" - ")
		sb.WriteString(e.Error())
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (e ConnErrors) Unwrap() []error {
	if len(e) == 0 {
		return nil
	}
	out := make([]error, len(e))
	for i := range e {
		// take the address so we have a *ConnError which implements error
		out[i] = &e[i]
	}
	return out
}

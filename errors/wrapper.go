package errors

import (
	"fmt"
)

// Wrapper provides a convenient mechanism for wrapping a number of erros with
// a common prefix.
type Wrapper string

// AppendedWith gives another Wrapper with the additional string appended to
// the prefix which will be added to any subsequent errors.
func (w Wrapper) AppendedWith(s string) Wrapper {
	return Wrapper(fmt.Sprintf("%s: %s", string(w), s))
}

// Wrap prepends the previously set prefix to the given non-nil error. It will
// refrain from wrapping a nil error.
func (w Wrapper) Wrap(err error) error {
	if nil != err {
		return fmt.Errorf("%s: %s", string(w), err.Error())
	}
	return nil
}

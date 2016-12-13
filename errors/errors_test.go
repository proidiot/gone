package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const error1 = New(
	"error1",
)

const error2 = ErrorString(
	"error2",
)

func errorIdentity(e error) error {
	return e
}

func TestErrors(t *testing.T) {
	error3 := New("error3")
	testData := map[string]error{
		"error1": error1,
		"error2": error2,
		"error3": error3,
		"error4": New("error4"),
		"error5": ErrorString("error5"),
		"error6": New(fmt.Sprintf("%s", "error6")),
	}

	for k, v := range testData {
		e := errorIdentity(v)
		assert.Error(t, e)
		assert.Equal(t, e, v)
		assert.Equal(t, e.Error(), k)
	}
}

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		"error6": New(fmt.Sprintf("error%d", 6)),
	}

	for k, v := range testData {
		e := errorIdentity(v)
		assert.Errorf(
			t,
			e,
			"Error creation test produced a non-error for: %s",
			k,
		)
		assert.Equal(
			t,
			e,
			v,
			"Error creation test produced an error that was not a"+
				" literal cast for: %s",
			k,
		)
		assert.Equal(
			t,
			e.Error(),
			k,
			"Error creation test produced an error with the wrong"+
				" string form for: %s",
			k,
		)
	}
}

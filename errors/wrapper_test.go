package errors

import (
	"fmt"
)

func ExampleWrapper() {
	wrapper := Wrapper("Error in wrapper test")
	appendedWrapper := wrapper.AppendedWith("Appended error")

	baseErr1 := New("Some error")
	err1 := appendedWrapper.Wrap(baseErr1)

	err2 := wrapper.Wrap(nil)

	baseErr3 := fmt.Errorf("Yet another %s", "error")
	err3 := wrapper.Wrap(baseErr3)

	fmt.Printf(
		"err1 is = %v\nerr2 is = %v\nerr3 is = %v\n",
		err1,
		err2,
		err3,
	)
	// Output:
	// err1 is = Error in wrapper test: Appended error: Some error
	// err2 is = <nil>
	// err3 is = Error in wrapper test: Yet another error
}

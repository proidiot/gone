package errors

// ErrorString is a minimal implementation of the error interface which can be
// used as either a variable or a constant.
//
// Since the original errors.errorString type was a struct and unexported, it
// is impossible to use it to create an error constant. However, since even the
// original errorString type stored only the string in question in its struct,
// nothing should be lost by literally defining the type as another name for a
// string. By declaring this ErrorString as an exported type defined to be a
// string, it is possible to use ErrorString as a variable or a constant.
type ErrorString string

// Error allows ErrorString to implement the error interface.
//
// Since a type needs only to implement Error and return a string, it is easy
// enough to make even a constant string trivially implement this type by
// having the function simply return its backing data structure re-cast as a
// string.
func (e ErrorString) Error() string {
	return string(e)
}

// New is a minimal implementation of the error interface which can be used as
// either a variable or a constant without many changes to existing code
// calling the similarly named constructor from the original errors package.
//
// Since the standard practice is to declare error variables using errors.New,
// a cheap hack to make this continue to work alongside the desired error
// constants given by ErrorString is to also declare New as an exported type
// defined to be a string. Even though this behavior is an arguably more
// significant departure from the original than the ErrorString behavior is,
// the syntax for declaring such an error should be identical, and the
// flexibility achieved thanks to Go's approach to interfaces should make this
// significant change in nature invisible to all other code.
type New string

// Error allows (the data type) New to implement the error interface.
//
// Just as with ErrorString, by implementing this function, this version of
// errors.New will look just as much like an error as the original despite a
// significant change in the nature of the underlying data type.
func (e New) Error() string {
	return string(e)
}

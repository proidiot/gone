// Package errors is intended to basically replace the default errors package.
// It is intended that one could replace an error import with an import of this
// package instead and code should (hopefully) not have to change as a result.
// In exchange, this slight tweak allows for error constants instead of just
// package-level error variables. It seems that doing so would likely result in
// much more noise-making in the event that a bug leads to code attempting to
// overwrite one of these errors.
package errors

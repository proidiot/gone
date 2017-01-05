// Package log provides a logging mechanism designed to have a similar feel to
// the original syslog calls in C, but with a few extra behaviors that could
// come in handy. The motivation for this piece is the notion that not only is
// it vital to make logging as pain-free as possible from the perspective of a
// developer, but to also defer as many decisions as possible about retention
// and routing to administrators.
package log

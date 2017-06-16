package syslogger

import (
	"github.com/proidiot/gone/log/pri"
)

// Syslogger allows log messages to be recorded (or forwarded to another
// Syslogger) in some way.
type Syslogger interface {
	Syslog(p pri.Priority, msg interface{}) error
}

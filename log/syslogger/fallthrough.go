package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

// Fallthrough is a syslogger.Syslogger that logs to a backup syslogger.Sylogger
// in the event that there is trouble logging to the default
// syslogger.Syslogger.
type Fallthrough struct {
	Default     Syslogger
	Fallthrough Syslogger
}

// Syslog logs a message. In the case of a Fallthrough, an attempt will be made
// to log to a default syslogger.Syslogger, then if that fails an attempt will
// be made to log to a fallthrough syslogger.Sylogger, then if that fails an
// error is returned.
func (f *Fallthrough) Syslog(p pri.Priority, msg interface{}) error {
	if f.Default != nil && f.Default.Syslog(p, msg) == nil {
		return nil
	} else if f.Fallthrough != nil {
		return f.Fallthrough.Syslog(p, msg)
	} else {
		return errors.New(
			"A syslogger.Fallthrough must have a non-nil" +
				" fallthrough syslogger in order to be" +
				" meaningful, but an attempt has been made to" +
				" write a log to a syslogger.Fallthrough with" +
				" at least a nil fallthrough syslogger.",
		)
	}
}

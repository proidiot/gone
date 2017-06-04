package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

type Fallthrough struct {
	Default     Syslogger
	Fallthrough Syslogger
}

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

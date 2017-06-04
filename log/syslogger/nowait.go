package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

type NoWait struct {
	Syslogger Syslogger
}

func (n *NoWait) Syslog(p pri.Priority, msg interface{}) error {
	if n.Syslogger != nil {
		go n.Syslogger.Syslog(p, msg)
		return nil
	} else {
		return errors.New(
			"A syslogger.NoWait must have a non-nil syslogger in" +
				" order to be meaningful, but an attempt has" +
				" been made to write a log to a" +
				" syslogger.NoWait with a nil syslogger.",
		)
	}
}

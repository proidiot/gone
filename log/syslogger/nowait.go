package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

// NoWait is a syslogger.Syslogger that allows calls to Syslog to return
// immediately.
type NoWait struct {
	Syslogger Syslogger
}

// Syslog logs a message. In the case of NoWait, the message will be sent to
// another syslogger.Syslogger asynchronously.
func (n *NoWait) Syslog(p pri.Priority, msg interface{}) error {
	if n.Syslogger == nil {
		return errors.New(
			"A syslogger.NoWait must have a non-nil syslogger in" +
				" order to be meaningful, but an attempt has" +
				" been made to write a log to a" +
				" syslogger.NoWait with a nil syslogger.",
		)
	}

	go n.Syslogger.Syslog(p, msg)
	return nil
}

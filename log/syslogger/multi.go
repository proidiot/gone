package syslogger

import (
	"github.com/proidiot/gone/log/pri"
)

// Multi is a syslogger.Syslogger that will send messages to all of several
// other syslogger.Sysloggers.
type Multi struct {
	Sysloggers []Syslogger
	TryAll     bool
}

// Syslog logs a message. In the case of Multi, the message will be sent to each
// of several other syslogger.Sysloggers assuming none of the
// syslogger.Sysloggers produce an error.
func (m *Multi) Syslog(p pri.Priority, msg interface{}) error {
	var err error

	for _, s := range m.Sysloggers {
		if e := s.Syslog(p, msg); e != nil {
			if !m.TryAll {
				return e
			} else if err == nil {
				err = e
			}
		}
	}

	return err
}

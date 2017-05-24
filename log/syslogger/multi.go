package syslogger

import (
	"github.com/proidiot/gone/log/pri"
)

type Multi struct {
	Sysloggers []Syslogger
}

func (m Multi) Syslog(p pri.Priority, msg interface{}) error {
	for _, s := range m.Sysloggers {
		if e := s.Syslog(p, msg); e != nil {
			return e
		}
	}

	return nil
}

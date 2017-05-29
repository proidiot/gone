package syslogger

import (
	"github.com/proidiot/gone/log/pri"
)

type Multi struct {
	Sysloggers []Syslogger
	TryAll     bool
}

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

package syslogger

import (
	"github.com/proidiot/gone/log/pri"
)

type Fallthrough struct {
	Default     Syslogger
	Fallthrough Syslogger
}

func (f *Fallthrough) Syslog(p pri.Priority, msg interface{}) error {
	if f.Default.Syslog(p, msg) == nil {
		return nil
	} else {
		return f.Fallthrough.Syslog(p, msg)
	}
}

package syslogger

import (
	"github.com/proidiot/gone/log/pri"
)

type NoWait struct {
	Syslogger Syslogger
}

func (n NoWait) Syslog(p pri.Priority, msg interface{}) error {
	go n.Syslogger.Syslog(p, msg)
	return nil
}

package syslogger

import (
	"github.com/proidiot/gone/log/pri"
	"os"
)

type Stderr struct {
}

func (s *Stderr) Syslog(p pri.Priority, msg interface{}) error {
	return (&Writer{
		os.Stderr,
	}).Syslog(p, msg)
}

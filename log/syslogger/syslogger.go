package syslogger

import (
	"github.com/proidiot/gone/log/pri"
)

type Syslogger interface {
	Syslog(p pri.Priority, msg interface{}) error
}

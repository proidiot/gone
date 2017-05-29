package syslogger

import (
	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/pri"
)

type SeverityMask struct {
	Syslogger Syslogger
	Mask      mask.Mask
}

func (s *SeverityMask) Syslog(p pri.Priority, msg interface{}) error {
	if s.Mask.Masked(p.Severity) {
		return nil
	}

	return s.Syslogger.Syslog(p, msg)
}

package syslogger

import (
	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/pri"
)

// SeverityMask is a syslogger.Syslogger that only forwards messages that aren't
// masked to another syslogger.Syslogger.
type SeverityMask struct {
	Syslogger Syslogger
	Mask      mask.Mask
}

// Syslog logs a message. In the case of SeverityMask, the message is sent to
// another syslogger.Syslogger if and only if the message isn't masked.
func (s *SeverityMask) Syslog(p pri.Priority, msg interface{}) error {
	if s.Mask.Masked(p.Severity()) {
		return nil
	}

	return s.Syslogger.Syslog(p, msg)
}

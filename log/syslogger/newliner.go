package syslogger

import (
	"fmt"
	"strings"

	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

// Newliner is a syslogger.Syslogger that assures the last byte in a message is
// a newline (i.e. a literal byte 0x0A).
type Newliner struct {
	Syslogger Syslogger
}

// Syslog logs a message. In the case of Newliner, the message has a newline
// added if not already present before passing through to another Syslogger.
func (n *Newliner) Syslog(p pri.Priority, msg interface{}) error {
	var s string
	switch m := msg.(type) {
	case fmt.Stringer:
		s = m.String()
	case string:
		s = m
	default:
		return errors.New(
			"The *syslogger.Newliner does not support message"+
				" types other than fmt.Stringer and string,"+
				" but the given message has a different type.",
		)
	}

	if strings.HasSuffix(s, "\n") {
		return n.Syslogger.Syslog(p, s)
	}

	return n.Syslogger.Syslog(p, s + "\n")
}

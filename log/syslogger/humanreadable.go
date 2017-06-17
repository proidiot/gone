package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"os"
	"time"
)

// HumanReadable is a syslogger.Syslogger that will format the message in a
// human readable way before passing the modified message to another
// syslogger.Syslogger.
type HumanReadable struct {
	Syslogger Syslogger
	Ident     string
	Facility  pri.Priority
	Pid       bool
}

// Syslog logs a message. In the case of HumanReadable, the message will be
// given a specific format and then forwarded to another syslogger.Syslogger.
func (h *HumanReadable) Syslog(p pri.Priority, msg interface{}) error {
	var s string
	switch msg := msg.(type) {
	case string:
		s = msg
	case fmt.Stringer:
		s = msg.String()
	case error:
		s = msg.Error()
	default:
		return errors.New(
			"The *syslogger.HumanReadable expects the message" +
				" argument to have the type string," +
				" fmt.Stringer, or error, but the given" +
				" message argument does not have one of" +
				" these types.",
		)
	}

	if p.ValidFacility() != nil || p.Facility() == 0x00 {
		if h.Facility == 0x00 {
			p = pri.User | p.Severity()
		} else {
			p = h.Facility | p.Severity()
		}
	}

	timestamp := time.Now().Format(time.UnixDate)

	hostname, e := osHostname()
	if e != nil {
		hostname = "localhost"
	}

	ident := h.Ident
	if ident == "" {
		ident = os.Args[0]
	}

	if h.Pid {
		ident = fmt.Sprintf(
			"%s[%d]",
			ident,
			os.Getpid(),
		)
	}

	m := fmt.Sprintf(
		"%s %s %s %s %s %s",
		p.Facility(),
		p.Severity(),
		timestamp,
		hostname,
		ident,
		s,
	)

	return h.Syslogger.Syslog(pri.Priority(0x0), m)
}

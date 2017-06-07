package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"os"
	"time"
)

var osHostname = os.Hostname

type HumanReadable struct {
	Syslogger Syslogger
	Facility  pri.Priority
	Ident     string
	Pid       bool
}

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

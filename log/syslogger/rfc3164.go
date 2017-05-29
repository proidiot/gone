package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"os"
	"time"
)

type Rfc3164 struct {
	Syslogger Syslogger
	Facility  pri.Facility
	Ident     string
	Pid       bool
}

func (r *Rfc3164) Syslog(p pri.Priority, msg interface{}) error {
	content, goodType := msg.(string)
	if !goodType {
		return errors.New(
			"The syslogger.Rfc3164 expects the message argument" +
				" to be a string, but the given message does" +
				" not have the string type.",
		)
	}

	if p.Facility.Valid() != nil || p.Facility == 0x00 {
		if r.Facility == 0x00 {
			p.Facility = pri.User
		} else {
			p.Facility = r.Facility
		}
	}
	prival, e := p.Combine()
	if e != nil {
		return e
	}

	timestamp := time.Now().Format(time.Stamp)

	hostname, e := os.Hostname()
	if e != nil {
		hostname = "localhost"
	}

	tag := r.Ident
	if tag == "" {
		tag = os.Args[0]
	}

	pid := ""
	if r.Pid {
		pid = fmt.Sprintf(
			"[%d]",
			os.Getpid(),
		)
	}

	res := fmt.Sprintf(
		"<%d>%s %s %s%s: %s",
		prival,
		timestamp,
		hostname,
		tag,
		pid,
		content,
	)

	if l := len([]byte(res)); l > 1024 {
		return errors.New(
			fmt.Sprintf(
				"The maximum total length of an RFC3164"+
					" syslog message is 1024 bytes, but"+
					" the generated syslog message has"+
					" total length %d bytes.",
				l,
			),
		)
	}

	return r.Syslogger.Syslog(pri.Priority{}, res)
}

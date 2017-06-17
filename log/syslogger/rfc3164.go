package syslogger

import (
	"fmt"
	"os"
	"time"

	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

// Rfc3164 is a syslogger.Syslogger that will format the message in a way that
// is intended to be compliant with RFC 3164 before passing the modified message
// to another syslogger.Syslogger.
type Rfc3164 struct {
	Syslogger Syslogger
	Ident     string
	Facility  pri.Priority
	Pid       bool
}

// Syslog logs a message. In the case of Rfc3164, the message is will be given a
// specific format and then forwarded to another syslogger.Syslogger.
func (r *Rfc3164) Syslog(p pri.Priority, msg interface{}) error {
	content, goodType := msg.(string)
	if !goodType {
		return errors.New(
			"The syslogger.Rfc3164 expects the message argument" +
				" to be a string, but the given message does" +
				" not have the string type.",
		)
	}

	if p.ValidFacility() != nil || p.Facility() == 0x00 {
		if r.Facility == 0x00 {
			p = pri.User | p.Severity()
		} else {
			p = r.Facility | p.Severity()
		}
	}

	timestamp := time.Now().Format(time.Stamp)

	hostname, e := osHostname()
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
		p,
		timestamp,
		hostname,
		tag,
		pid,
		content,
	)

	if l := len([]byte(res)); l > 1024 {
		return fmt.Errorf(
			"The maximum total length of an RFC3164 syslog"+
				" message is 1024 bytes, but the generated"+
				" syslog message has total length %d bytes.",
			l,
		)
	}

	return r.Syslogger.Syslog(pri.Priority(0x0), res)
}

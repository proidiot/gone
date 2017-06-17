package syslogger

import (
	"io"

	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

// Writer is a syslogger.Syslogger that writes messages directly to an
// io.Writer.
type Writer struct {
	Writer io.Writer
}

// Syslog logs a message. In the case of Writer, the message is written directly
// to an arbitrary io.Writer.
func (w *Writer) Syslog(p pri.Priority, msg interface{}) error {
	if p != 0x00 {
		return errors.New(
			"The basic syslog.Writer cannot differentiate" +
				" between log priorities so it expects a" +
				" zero-valued priority argument, but a" +
				" non-zero pri.Priority was given.",
		)
	}

	switch m := msg.(type) {
	case string:
		_, e := io.WriteString(w.Writer, m)
		return e
	case []byte:
		_, e := w.Writer.Write(m)
		return e
	default:
		return errors.New(
			"The basic *syslogger.Writer does not support" +
				" message typesother than string and []byte," +
				" but the given message has a different type.",
		)
	}
}

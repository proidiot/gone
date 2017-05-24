package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"io"
)

type Writer struct {
	Writer io.Writer
}

func (w Writer) Syslog(p pri.Priority, msg interface{}) error {
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
			" message typesother than string and []byte, but the" +
			" given message has a different type.",
		)
	}
}

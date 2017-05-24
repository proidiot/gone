package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"io"
	"os"
)

type DevConsole struct {
	w Writer
}

func (c DevConsole) Syslog(p pri.Priority, msg interface{}) error {
	if c.w.Writer == nil {
		return errors.New(
			"The *syslogger.DevConsole must be initialized" +
			" before use.",
		)
	}

	return c.w.Syslog(p, msg)
}

func (c DevConsole) Close() error {
	return c.w.Writer.(io.Closer).Close()
}

func NewDevConsole() (DevConsole, error) {
	f, e := os.OpenFile("/dev/console", os.O_APPEND | os.O_WRONLY, 0600)
	if e != nil {
		return DevConsole{}, e
	}

	return DevConsole{Writer{f}}, nil
}

package syslogger

import (
	"github.com/proidiot/gone/log/pri"
	"io"
)

type WriteCloser struct {
	WriteCloser io.WriteCloser
}

func (w *WriteCloser) Syslog(p pri.Priority, msg interface{}) error {
	return (&Writer{w.WriteCloser}).Syslog(p, msg)
}

func (w *WriteCloser) Close() error {
	return w.WriteCloser.Close()
}

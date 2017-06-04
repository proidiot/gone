package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

type flagSyslog struct {
	flag bool
}

func (f *flagSyslog) Syslog(p pri.Priority, msg interface{}) error {
	f.flag = true
	return nil
}

type errorSyslog struct {
}

func (e errorSyslog) Syslog(p pri.Priority, msg interface{}) error {
	return errors.New("Syslog called on an errorSyslog")
}

type errorWriter struct {
}

func (e errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("Writing to an errorWriter")
}

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

type syncFlagSyslog struct {
	flag bool
	sync chan interface{}
}

func (s *syncFlagSyslog) Syslog(p pri.Priority, msg interface{}) error {
	<-s.sync
	s.flag = true
	return nil
}

func newSyncFlagSyslog() *syncFlagSyslog {
	return &syncFlagSyslog{
		sync: make(chan interface{}),
	}
}

type errorSyslogError struct {
	e errors.New
	s *errorSyslog
}

func (e *errorSyslogError) Error() string {
	return e.e.Error()
}

type errorSyslog struct {
}

func (es *errorSyslog) Syslog(p pri.Priority, msg interface{}) error {
	return &errorSyslogError{
		e: "Syslog called on an errorSyslog",
		s: es,
	}
}

type errorWriter struct {
}

func (e errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("Writing to an errorWriter")
}

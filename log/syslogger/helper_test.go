package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
)

type flagSyslogger struct {
	Flag bool
}

func (f *flagSyslogger) Syslog(p pri.Priority, msg interface{}) error {
	f.Flag = true
	return nil
}

type syncFlagSyslogger struct {
	Flag bool
	sync <-chan interface{}
}

func (s *syncFlagSyslogger) Syslog(p pri.Priority, msg interface{}) error {
	<-s.sync
	s.Flag = true
	return nil
}

func newSyncFlagSyslogger() (*syncFlagSyslogger, chan<- interface{}) {
	sc := make(chan interface{})
	return &syncFlagSyslogger{
		sync: sc,
	}, sc
}

type syncCountSyslogger struct {
	Count uint
	sync  <-chan interface{}
}

func (s *syncCountSyslogger) Syslog(p pri.Priority, msg interface{}) error {
	<-s.sync
	s.Count++
	return nil
}

func newSyncCountSyslogger() (*syncCountSyslogger, chan<- interface{}) {
	sc := make(chan interface{})
	return &syncCountSyslogger{
		sync: sc,
	}, sc
}

type errorSysloggerError struct {
	E errors.New
	S *errorSyslogger
}

func (e *errorSysloggerError) Error() string {
	return e.E.Error()
}

type errorSyslogger struct {
}

func (es *errorSyslogger) Syslog(p pri.Priority, msg interface{}) error {
	return &errorSysloggerError{
		E: "Syslog called on an errorSyslog",
		S: es,
	}
}

type errorWriter struct {
}

func (e errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("Writing to an errorWriter")
}

type recordStringSyslogger struct {
	P pri.Priority
	M string
}

func (rs *recordStringSyslogger) Syslog(p pri.Priority, msg interface{}) error {
	s, ok := msg.(string)
	if !ok {
		return errors.New("Non-string passed to a recordStringSyslog")
	}

	rs.M = s
	rs.P = p
	return nil
}

type errorCloser struct {
}

func (e *errorCloser) Close() error {
	return errors.New("Closing an errorCloser")
}

type stringer struct {
	S string
}

func (s *stringer) String() string {
	return string(s.S)
}

package log

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/opt"
	"github.com/proidiot/gone/log/pri"
)

type testSyslogger struct {
	LastPri      pri.Priority
	LastMsg      interface{}
	TriggerError bool
}

func (t *testSyslogger) triggerError() error {
	if t.TriggerError {
		return errors.New("Artificial error triggered in testSyslogger")
	}

	return nil
}

func (t *testSyslogger) Syslog(p pri.Priority, msg interface{}) error {
	t.LastPri = p
	t.LastMsg = msg
	return t.triggerError()
}

func (t *testSyslogger) Openlog(string, opt.Option, pri.Priority) error {
	return t.triggerError()
}

func (t *testSyslogger) Close() error {
	return t.triggerError()
}

type limitedSyslogger struct {
	LastPri      pri.Priority
	LastMsg      interface{}
	TriggerError bool
}

func (l *limitedSyslogger) triggerError() error {
	if l.TriggerError {
		return errors.New(
			"Artificial error triggered in limitedSyslogger",
		)
	}

	return nil
}

func (l *limitedSyslogger) Syslog(p pri.Priority, msg interface{}) error {
	l.LastPri = p
	l.LastMsg = msg
	return l.triggerError()
}

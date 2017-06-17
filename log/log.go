package log

import (
	"os"

	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/opt"
	"github.com/proidiot/gone/log/pri"
	"github.com/proidiot/gone/log/syslogger"
)

var log syslogger.Syslogger

func init() {
	l := &syslogger.Posixish{}
	e := l.Openlog(
		os.Getenv("LOG_IDENT"),
		opt.GetFromEnv(),
		pri.GetFromEnv(),
	)
	if e != nil {
		panic(e)
	}
	e = l.SetLogMask(mask.GetFromEnv())
	if e != nil {
		panic(e)
	}

	log = l
}

// SetSyslogger overwrites the default global syslogger.Syslogger with the one
// given explicitly.
func SetSyslogger(s syslogger.Syslogger) {
	log = s
}

// Openlog allows the global syslogger.Syslogger to be reset with certain
// explicit initialization values.
func Openlog(ident string, o opt.Option, f pri.Priority) error {
	type openlogger interface {
		Openlog(string, opt.Option, pri.Priority) error
	}

	ol, ok := log.(openlogger)
	if !ok {
		return errors.New(
			"Default global log has been set to a" +
				" syslogger.Syslogger without an Openlog" +
				" function.",
		)
	}

	return ol.Openlog(ident, o, f)
}

// Syslog allows logs to be written to the global syslogger.Syslogger.
func Syslog(p pri.Priority, msg interface{}) error {
	return log.Syslog(p, msg)
}

// Closelog ends the log session of the global syslogger.Syslogger. Depending on
// which kind of syslogger.Syslogger the default is set to, this could result in
// all future Syslog calls creating errors, or it could have practically no
// effect.
func Closelog() error {
	type closer interface {
		Close() error
	}

	cl, ok := log.(closer)
	if !ok {
		return errors.New(
			"Default global log has been set to a" +
				" syslogger.Syslogger without a Close" +
				" function.",
		)
	}

	return cl.Close()
}

// Emerg sends a log message with priority Emerg
func Emerg(m interface{}) error {
	return Syslog(pri.Emerg, m)
}

// Emergency sends a log message with priority Emerg
func Emergency(m interface{}) error {
	return Emerg(m)
}

// Alert sends a log message with priority Alert
func Alert(m interface{}) error {
	return Syslog(pri.Alert, m)
}

// Crit sends a log message with priority Crit
func Crit(m interface{}) error {
	return Syslog(pri.Crit, m)
}

// Critical sends a log message with priority Crit
func Critical(m interface{}) error {
	return Crit(m)
}

// Err sends a log message with priority Err
func Err(m interface{}) error {
	return Syslog(pri.Err, m)
}

// Error sends a log message with priority Err
func Error(m interface{}) error {
	return Err(m)
}

// Warning sends a log message with priority Warning
func Warning(m interface{}) error {
	return Syslog(pri.Warning, m)
}

// Warn sends a log message with priority Warning
func Warn(m interface{}) error {
	return Warning(m)
}

// Notice sends a log message with priority Notice
func Notice(m interface{}) error {
	return Syslog(pri.Notice, m)
}

// Info sends a log message with priority Info
func Info(m interface{}) error {
	return Syslog(pri.Info, m)
}

// Information sends a log message with priority Info
func Information(m interface{}) error {
	return Info(m)
}

// Debug sends a log message with priority Debug
func Debug(m interface{}) error {
	return Syslog(pri.Debug, m)
}

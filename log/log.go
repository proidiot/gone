package log

import (
	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/opt"
	"github.com/proidiot/gone/log/pri"
	"github.com/proidiot/gone/log/syslogger"
)

var log syslogger.Syslogger

func init() {
	log = syslogger.Posixish{}
	e := log.Openlog(
		os.Getenv("LOG_IDENT"),
		opt.GetFromEnv(),
		pri.GetFromEnv(),
	)
	if e != nil {
		panic(e)
	}
	e = log.SetLogMask(mask.GetFromEnv())
	if e != nil {
		panic(e)
	}
}

func Openlog(ident string, o opt.Option, f pri.Facility) error {
	return log.Openlog(ident, o, f)
}

func Syslog(p pri.Priority, msg interface{}) error {
	return log.Syslog(p, msg)
}

func Closelog() error {
	l, goodType := log.(io.Closer)
	if goodType {
		return l.Close()
	} else {
		// TODO
		return nil
	}
}

func Emerg(m interface{}) error {
	return Syslog(pri.Emerg, m)
}

func Emergency(m interface{}) error {
	return Emerg(m)
}

func Alert(m interface{}) error {
	return Syslog(pri.Alert, m)
}

func Crit(m interface{}) error {
	return Syslog(pri.Crit, m)
}

func Critical(m interface{}) error {
	return Crit(m)
}

func Err(m interface{}) error {
	return Syslog(pri.Err, m)
}

func Error(m interface{}) error {
	return Err(m)
}

func Warning(m interface{}) error {
	return Syslog(pri.Warning, m)
}

func Warn(m interface{}) error {
	return Warning(m)
}

func Notice(m interface{}) error {
	return Syslog(pri.Notice, m)
}

func Info(m interface{}) error {
	return Syslog(pri.Info, m)
}

func Information(m interface{}) error {
	return Info(m)
}

func Debug(m interface{}) error {
	return Syslog(pri.Debug, m)
}

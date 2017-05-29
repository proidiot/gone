package log

import (
	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/opt"
	"github.com/proidiot/gone/log/pri"
	"github.com/proidiot/gone/log/syslogger"
	"os"
)

var log syslogger.Posixish

func init() {
	l := syslogger.Posixish{}
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

func Openlog(ident string, o opt.Option, f pri.Facility) error {
	return log.Openlog(ident, o, f)
}

func Syslog(p pri.Priority, msg interface{}) error {
	return log.Syslog(p, msg)
}

func Closelog() error {
	return log.Close()
}

func Emerg(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Emerg}, m)
}

func Emergency(m interface{}) error {
	return Emerg(m)
}

func Alert(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Alert}, m)
}

func Crit(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Crit}, m)
}

func Critical(m interface{}) error {
	return Crit(m)
}

func Err(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Err}, m)
}

func Error(m interface{}) error {
	return Err(m)
}

func Warning(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Warning}, m)
}

func Warn(m interface{}) error {
	return Warning(m)
}

func Notice(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Notice}, m)
}

func Info(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Info}, m)
}

func Information(m interface{}) error {
	return Info(m)
}

func Debug(m interface{}) error {
	return Syslog(pri.Priority{0x00,pri.Debug}, m)
}

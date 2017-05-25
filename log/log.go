package log

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"log/syslog"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	// local options
	defaultFacility = syslog.LOG_USER // Default syslog facility in C.
)

type logData struct {
	p byte
	m string
}

type Settings struct {
	Dial        bool     // affects writer
	Network     string   // affects writer
	Raddr       string   // affects writer
	Appname     string   // affects logger
	Facility    Facility // affects logger
	Ident       string   // affects logger
	Pid         bool     // affects logger
	Cons        bool     // affects writer
	NoDelay     bool     // affects writer
	NoWait      bool     // ?
	Perror      bool     // affects writer
	NoFallback  bool     // affects writer
	Reestablish bool     // affects writer
}

type syslogger struct {
	c       chan *logData
	s       Settings
	w       *syslog.Writer
	l       *sync.RWMutex
	console *os.File
	appname string
}

var log syslogger

// Stored Openlog arguments (in case of delay flag)
var savedIdent = os.Args[0]         // C default
var savedFacility = syslog.LOG_USER // C default

// Openlog
func Openlog(ident string, option Option, facility Facility) error {
	if (option&LOG_ODELAY) != 0 && (option&LOG_NDELAY) != 0 {
		e := errors.New(
			"LOG_ODELAY and LOG_NDELAY are both being passed to" +
				"Openlog, but these options are mutually exclusive.",
		)

		Syslog(LOG_ERR, e.Error())

		return e
	} else if facility > (LOG_LOCAL7 + Facility(LOG_DEBUG)) {
		e := errors.New(
			"Openlog expects the facility argument to be one of" +
				" the 24 syslog facility constants, but some other" +
				" value has been specified fo a call to Openlog.",
		)

		Syslog(LOG_ERR, e.Error())

		return e
	} else if (facility & Facility(SeverityMask)) != 0 {
		e := errors.New(
			"Openlog does not accept a combined facility and" +
				" severity value for its facility argument even" +
				" though Syslog does accept such a combined value" +
				" for its priority argument, but such a combined" +
				" value has been given for a call to Openlog.",
		)

		Syslog(LOG_ERR, e.Error())

		return e
	} else if (option & OptionsMask) != option {
		e := errors.New(
			"Openlog expects its option argument to be the" +
				" bitwise OR of known option flags, but other bits" +
				" were set in the option argument for a call to" +
				" Openlog.",
		)

		Syslog(LOG_ERR, e.Error())

		return e
	}

	s := Settings{
		Facility:    byte(facility & FacilityMask),
		Ident:       ident,
		Pid:         (option & LOG_PID) != 0,
		Cons:        (option & LOG_CONS) != 0,
		NoDelay:     (option & LOG_NDELAY) != 0,
		NoWait:      (option & LOG_NOWAIT) != 0,
		Perror:      (option & LOG_PERROR) != 0,
		NoFallback:  (option & LOG_NOFALLBACK) != 0,
		Reestablish: (option & LOG_REESTABLISH) != 0,
	}

	return SetSettings(s)
}

func SetSettings(s Settings) error {
	log.l.Lock()
	defer log.l.Unlock()

	if s.NoDelay {
		if e := establishLogger(s); e != nil {
			return e
		}
	}

	log.s = s

	return nil
}

func establishLogger(s Settings) error {
	var e error
	if s.Dial {
		w, e := syslog.Dial(
			s.Network,
			s.Raddr,
			syslog.Priority(int(s.Facility)),
			s.Ident,
		)

		if e == nil {
			log.w = w

			return nil
		}

	} else {
		w, e := syslog.New(
			syslog.Priority(int(s.Facility)),
			s.Ident,
		)

		if e == nil {
			log.w = w
			return nil
		}
	}

	if s.Cons {
		f, e := os.OpenFile(
			"/dev/console",
			os.O_APPEND|os.O_WRONLY,
			0600,
		)

		if e == nil {
			log.console = f

			return nil
		}
	}

	if s.NoFallback {
		return e
	} else {
		return nil
	}
}

func (q *syslogger) syslogFmt(priority Priority, m string) (string, error) {
	prival := priority & Priority(FacilityMask)
	if prival == 0 || prival > Priority(LOG_LOCAL7) {
		prival = Priority(l.s.Facility)
	}
	prival = prival | (priority & SeverityMask)

	// BUG fix
	timestamp := "-"

	hostname, e := os.Hostname()
	if e != nil {
		hostname = "-"
	}

	appname := l.appname

	procid := "-"
	if l.s.Pid {
		procid := strconv.Itoa(os.Getpid())
	}

	msgid := l.s.Ident

	// Maybe add this later?
	structureddata := "-"

	msg := m

	// See RFC 5424.
	res := fmt.Sprintf(
		"<%d>1 %s %s %s %s %s %s %s",
		prival,
		timestamp,
		hostname,
		appname,
		procid,
		msgid,
		structureddata,
		msg,
	)

	return res, nil
}

func syslogg(priority Priority, m string) error {
	var e error // needed for switch scope, here so syslogger won't be local

	if syslogger == nil {
		syslogger, e = syslog.New(savedFacility, savedIdent)
	}

	if e != nil {
		if savedCons {
			os.Stderr.Write([]byte(e.Error()))
		}
	} else {
		priorityMask := 7 // all 0's except the least 3 bits
		switch priority & priorityMask {
		case LOG_EMERG:
			e = syslogger.Emerg(m)
		case LOG_ALERT:
			e = syslogger.Alert(m)
		case LOG_CRIT:
			e = syslogger.Crit(m)
		case LOG_ERR:
			e = syslogger.Err(m)
		case LOG_WARNING:
			e = syslogger.Warning(m)
		case LOG_NOTICE:
			e = syslogger.Notice(m)
		case LOG_INFO:
			e = syslogger.Info(m)
		case LOG_DEBUG:
			e = syslogger.Debug(m)
		}
	}

	if savedPerror {
		var ePerror error

		if savedPid {
			_, ePerror = os.Stderr.Write(
				[]byte(
					fmt.Sprintf(
						"<%d>%s %s[%d]: %s",
						priority,
						time.Now().Format(time.Stamp),
						savedIdent,
						os.Getpid(),
						m,
					),
				),
			)
		} else {
			_, ePerror = os.Stderr.Write(
				[]byte(
					fmt.Sprintf(
						"<%d>%s %s: %s",
						priority,
						time.Now().Format(time.Stamp),
						savedIdent,
						m,
					),
				),
			)
		}

		if e == nil {
			return ePerror
		}
	}

	return e
}

func Syslog(p Priority, m interface{}) error {
	var s string

	switch m := m.(type) {
	case string:
		s = m
	case fmt.Stringer:
		s = m.String()
	case error:
		s = m.Error()
	default:
		e := errors.New(
			fmt.Sprintf(
				"Syslog exoects its msg argument to be one"+
					" of a string, a fmt.Stringer, or an error,"+
					" but a value was given to Syslog that does"+
					" not match any of these types: %v",
				m,
			),
		)

		return e
	}

	return syslog(p, s)
}

func Closelog() error {
	return syslogger.Close()
}

func Emerg(m interface{}) error {
	return Syslog(LOG_EMERG, m)
}

func Emergency(m interface{}) error {
	return Emerg(m)
}

func Alert(m interface{}) error {
	return Syslog(LOG_ALERT, m)
}

func Crit(m interface{}) error {
	return Syslog(LOG_CRIT, m)
}

func Critical(m interface{}) error {
	return Crit(m)
}

func Err(m interface{}) error {
	return Syslog(LOG_ERR, m)
}

func Error(m interface{}) error {
	return Err(m)
}

func Warning(m interface{}) error {
	return Syslog(LOG_WARNING, m)
}

func Warn(m interface{}) error {
	return Warning(m)
}

func Notice(m interface{}) error {
	return Syslog(LOG_NOTICE, m)
}

func Info(m interface{}) error {
	return Syslog(LOG_INFO, m)
}

func Information(m interface{}) error {
	return Info(m)
}

func Debug(m interface{}) error {
	return Syslog(LOG_DEBUG, m)
}

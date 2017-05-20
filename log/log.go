package log

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"log/syslog"
	"os"
	"time"
)

type logData struct {
	pri int
	m string
}

type logSettings struct {
	pid bool
	cons bool
	delay bool

}

var log chan<- *logData
var logControl chan<- *logSettings
var syslogger *syslog.Writer

// Flags
var savedPid = false // C default for perror, but log/syslog always sends PID
var savedCons = false // C default
var savedPerror = false // C default

// Stored Openlog arguments (in case of delay flag)
var savedIdent = os.Args[0] // C default
var savedFacility = syslog.LOG_USER // C default

const (
	// Log levels
	LOG_EMERG = 0
	LOG_ALERT = 1
	LOG_CRIT = 2
	LOG_ERR = 3
	LOG_WARNING = 4
	LOG_NOTICE = 5
	LOG_INFO = 6
	LOG_DEBUG = 7

	// Log facilities
	LOG_KERN = (0<<3)
	LOG_USER = (1<<3)
	LOG_MAIL = (2<<3)
	LOG_DAEMON = (3<<3)
	LOG_AUTH = (4<<3)
	LOG_SYSLOG = (5<<3)
	LOG_LPR = (6<<3)
	LOG_NEWS = (7<<3)
	LOG_UUCP = (8<<3)
	LOG_CRON = (9<<3)
	LOG_AUTHPRIV = (10<<3)
	LOG_FTP = (11<<3)
	LOG_NTP = (12<<3)
	LOG_LOG_AUDIT = (13<<3)
	LOG_LOG_ALERT = (14<<3)
	LOG_CLOCKD = (15<<3)
	LOG_LOCAL0 = (16<<3)
	LOG_LOCAL1 = (17<<3)
	LOG_LOCAL2 = (18<<3)
	LOG_LOCAL3 = (19<<3)
	LOG_LOCAL4 = (20<<3)
	LOG_LOCAL5 = (21<<3)
	LOG_LOCAL6 = (22<<3)
	LOG_LOCAL7 = (23<<3)

	// Log options
	LOG_PID = 0x01
	LOG_CONS = 0x02
	LOG_ODELAY = 0x04
	LOG_NDELAY = 0x08
	LOG_NOWAIT = 0x10
	LOG_PERROR = 0x20

	// local options
	defaultFacility = syslog.LOG_USER // Default syslog facility in C.
)

// Openlog 
func Openlog(ident string, option, facility int) (error) {
	var e error

	savedPid = (option & LOG_PID) != 0
	savedCons = (option & LOG_CONS) != 0
	savedPerror = (option & LOG_PERROR) != 0

	if (option & LOG_ODELAY) != 0 && (option & LOG_NDELAY) != 0 {
		e = errors.New(
			"LOG_ODELAY and LOG_NDELAY are both passed to" +
			"Openlog, but these options are mutually exclusive.",
		)

		if savedCons {
			os.Stderr.Write([]byte(e.Error()))
		}
	}
	delay := (option & LOG_NDELAY) == 0

	savedIdent = ident

	// Unary bitwise xor in golang is bitwise complement,
	// so ^7 is all 1's except the least 3 bits.
	facilityMask := ^7

	// Unfortunately, log/syslog is missing some facility codes according
	// to RFC5424. As such, this ugly switch is needed.
	switch facility & facilityMask {
	case LOG_KERN:
		savedFacility = syslog.LOG_KERN
	case LOG_USER:
		savedFacility = syslog.LOG_USER
	case LOG_MAIL:
		savedFacility = syslog.LOG_MAIL
	case LOG_DAEMON:
		savedFacility = syslog.LOG_DAEMON
	case LOG_AUTH:
		savedFacility = syslog.LOG_AUTH
	case LOG_SYSLOG:
		savedFacility = syslog.LOG_SYSLOG
	case LOG_LPR:
		savedFacility = syslog.LOG_LPR
	case LOG_NEWS:
		savedFacility = syslog.LOG_NEWS
	case LOG_UUCP:
		savedFacility = syslog.LOG_UUCP
	case LOG_CRON:
		savedFacility = syslog.LOG_CRON
	case LOG_AUTHPRIV:
		savedFacility = syslog.LOG_AUTHPRIV
	case LOG_FTP:
		savedFacility = syslog.LOG_FTP
	case LOG_NTP:
		savedFacility = syslog.LOG_DAEMON // missing, least bad alt.
	case LOG_LOG_AUDIT:
		savedFacility = syslog.LOG_AUTHPRIV // missing, least bad alt.
	case LOG_LOG_ALERT:
		savedFacility = syslog.LOG_USER // missing, least bad alt?
	case LOG_CLOCKD:
		savedFacility = syslog.LOG_CRON // missing, least bad alt.
	case LOG_LOCAL0:
		savedFacility = syslog.LOG_LOCAL0
	case LOG_LOCAL1:
		savedFacility = syslog.LOG_LOCAL1
	case LOG_LOCAL2:
		savedFacility = syslog.LOG_LOCAL2
	case LOG_LOCAL3:
		savedFacility = syslog.LOG_LOCAL3
	case LOG_LOCAL4:
		savedFacility = syslog.LOG_LOCAL4
	case LOG_LOCAL5:
		savedFacility = syslog.LOG_LOCAL5
	case LOG_LOCAL6:
		savedFacility = syslog.LOG_LOCAL6
	case LOG_LOCAL7:
		savedFacility = syslog.LOG_LOCAL7
	default:
		savedFacility = syslog.LOG_USER // C default rather than error
	}

	if !delay {
		syslogger, e = syslog.New(savedFacility, savedIdent)

		if savedCons {
			os.Stderr.Write([]byte(e.Error()))
		}
	}

	return e
}

func Syslog(priority int, m string) (error) {
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

func Closelog() (error) {
	return syslogger.Close()
}

func Emerg(m string) (error) {
	return Syslog(LOG_EMERG, m)
}

func Emergency(m string) (error) {
	return Emerg(m)
}

func Alert(m string) (error) {
	return Syslog(LOG_ALERT, m)
}

func Crit(m string) (error) {
	return Syslog(LOG_CRIT, m)
}

func Critical(m string) (error) {
	return Crit(m)
}

func Err(m string) (error) {
	return Syslog(LOG_ERR, m)
}

func Error(m string) (error) {
	return Err(m)
}

func Warning(m string) (error) {
	return Syslog(LOG_WARNING, m)
}

func Warn(m string) (error) {
	return Warning(m)
}

func Notice(m string) (error) {
	return Syslog(LOG_NOTICE, m)
}

func Info(m string) (error) {
	return Syslog(LOG_INFO, m)
}

func Information(m string) (error) {
	return Info(m)
}

func Debug(m string) (error) {
	return Syslog(LOG_DEBUG, m)
}


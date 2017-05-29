package pri

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"os"
)

// Log Priority (bitwise or of facility and severity)
// See RFC 3164 Section 4.1.1, RFC 5424 Section 6.2.1, and RFC 5427 Section 3.
type Priority byte

const (
	// Emerg represents an emergency; system is unusable.
	Emerg Priority = 0
	// Alert represents than an action must be taken immediately.
	Alert Priority = 1
	// Crit represents a critical condition.
	Crit Priority = 2
	// Err represents an error condition.
	Err Priority = 3
	// Warning represents a warning condition.
	Warning Priority = 4
	// Notice represents a normal but significant condition.
	Notice Priority = 5
	// Info represents an informational message.
	Info Priority = 6
	// Debug represents debug-level messages.
	Debug Priority = 7
)

const (
	// Kern represents kernel messages.
	Kern Priority = (0 * 8)

	// User represents user-level messages.
	User Priority = (1 * 8)

	// Mail represents mail system messages.
	Mail Priority = (2 * 8)

	// Daemon represents system daemons' messages.
	Daemon Priority = (3 * 8)

	// Auth represents authorization messages.
	Auth Priority = (4 * 8)

	// Syslog represents messages generated internally by syslogd.
	Syslog Priority = (5 * 8)

	// Lpr represents line printer subsystem messages.
	Lpr Priority = (6 * 8)

	// News represents network news subsystem messages.
	News Priority = (7 * 8)

	// Uucp represents UUCP subsystem messages.
	Uucp Priority = (8 * 8)

	// Cron represents clock daemon messages.
	Cron Priority = (9 * 8)

	// Authpriv represents security/authorization messages.
	Authpriv Priority = (10 * 8)

	// Ftp represents ftp daemon messages.
	Ftp Priority = (11 * 8)

	// Ntp represents NTP subsystem messages.
	Ntp Priority = (12 * 8)

	// Audit represents audit messages.
	Audit Priority = (13 * 8)

	// Console represents console messages.
	Console Priority = (14 * 8)

	// Cron2 represents clock daemon messages.
	Cron2 Priority = (15 * 8)

	// Local0 represents messages designated as local use 0.
	Local0 Priority = (16 * 8)

	// Local1 represents messages designated as local use 1.
	Local1 Priority = (17 * 8)

	// Local2 represents messages designated as local use 2.
	Local2 Priority = (18 * 8)

	// Local3 represents messages designated as local use 3.
	Local3 Priority = (19 * 8)

	// Local4 represents messages designated as local use 4.
	Local4 Priority = (20 * 8)

	// Local5 represents messages designated as local use 5.
	Local5 Priority = (21 * 8)

	// Local6 represents messages designated as local use 6.
	Local6 Priority = (22 * 8)

	// Local7 represents messages designated as local use 7.
	Local7 Priority = (23 * 8)

	InvalidFacility = errors.New(
		"The facility portion of a pri.Priority can only have one of" +
			" the 24 values from pri.Kern through pri.Local7." +
			" The acceptable values are described at length by" +
			" both RFC3164 and RFC5424, and these constants are" +
			" listed at:" +
			" https://godoc.org/github.com/proidiot/gone/log/pri",
	)
)

var lookupFacility = map[Priority]string{
	Kern:     "LOG_KERN",
	User:     "LOG_USER",
	Mail:     "LOG_MAIL",
	Daemon:   "LOG_DAEMON",
	Auth:     "LOG_AUTH",
	Syslog:   "LOG_SYSLOG",
	Lpr:      "LOG_LPR",
	News:     "LOG_NEWS",
	Uucp:     "LOG_UUCP",
	Cron:     "LOG_CRON",
	Authpriv: "LOG_AUTHPRIV",
	Ftp:      "LOG_FTP",
	Ntp:      "LOG_NTP",
	Audit:    "LOG_AUDIT",
	Console:  "LOG_CONSOLE",
	Cron2:    "LOG_CRON2",
	Local0:   "LOG_LOCAL0",
	Local1:   "LOG_LOCAL1",
	Local2:   "LOG_LOCAL2",
	Local3:   "LOG_LOCAL3",
	Local4:   "LOG_LOCAL4",
	Local5:   "LOG_LOCAL5",
	Local6:   "LOG_LOCAL6",
	Local7:   "LOG_LOCAL7",
}

var lookupSeverity = map[Priority]string{
	Emerg:   "LOG_EMERG",
	Alert:   "LOG_ALERT",
	Crit:    "LOG_CRIT",
	Err:     "LOG_ERR",
	Warning: "LOG_WARNING",
	Notice:  "LOG_NOTICE",
	Info:    "LOG_INFO",
	Debug:   "LOG_DEBUG",
}

func (p Priority) Facility() Priority {
	// Equivalent to pri.Kern | pri.User | ... | pri.Local7
	mask := Priority(0xF8)
	return p & mask
}

func (p Priority) ValidFacility() error {
	if p != p.Facility() || p > Local7 {
		return InvalidFacility
	} else {
		return nil
	}
}

func (p Priority) Severity() Priority {
	// Equivalent to LOG_EMERG | LOG_ALERT | ... | LOG_DEBUG
	mask := Priority(0x07)
	return p & mask
}

func (p Priority) String() string {
	res := ""
	if p.Facility() != 0 {
		if fstring, present := lookupFacility[p.Facility()]; present {
			res += fstring
		} else {
			res += fmt.Sprintf("Priority(%#x)", byte(p.Facility()))
		}

		if p.Severity() != 0 {
			res += "|"
		} else {
			return res
		}
	}

	res += lookupSeverity[p.Severity()]

	return res
}

func GetFromEnv() Priority {
	if val, set := os.LookupEnv("LOG_FACILITY"); set {
		for p, pstring := range lookupFacility {
			if pstring == val {
				return p
			}
		}
	} else if val, set := os.LookupEnv("LOG_PRIORITY"); set {
		for p, pstring := range lookupFacility {
			if pstring == val {
				return p
			}
		}
	}

	return User
}

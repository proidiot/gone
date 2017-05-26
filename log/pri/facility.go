//go:generate stringer -type Facility
package pri

import (
	"github.com/proidiot/gone/errors"
	"os"
)

// Log facility
// See RFC 3164 Section 4.1.1, RFC 5424 Section 6.2.1, and RFC 5427 Section 3.
type Facility byte

const (
	// Kern represents kernel messages.
	Kern Facility = (0 * 8)

	// User represents user-level messages.
	User Facility = (1 * 8)

	// Mail represents mail system messages.
	Mail Facility = (2 * 8)

	// Daemon represents system daemons' messages.
	Daemon Facility = (3 * 8)

	// Auth represents authorization messages.
	Auth Facility = (4 * 8)

	// Syslog represents messages generated internally by syslogd.
	Syslog Facility = (5 * 8)

	// Lpr represents line printer subsystem messages.
	Lpr Facility = (6 * 8)

	// News represents network news subsystem messages.
	News Facility = (7 * 8)

	// Uucp represents UUCP subsystem messages.
	Uucp Facility = (8 * 8)

	// Cron represents clock daemon messages.
	Cron Facility = (9 * 8)

	// Authpriv represents security/authorization messages.
	Authpriv Facility = (10 * 8)

	// Ftp represents ftp daemon messages.
	Ftp Facility = (11 * 8)

	// Ntp represents NTP subsystem messages.
	Ntp Facility = (12 * 8)

	// Audit represents audit messages.
	Audit Facility = (13 * 8)

	// Console represents console messages.
	Console Facility = (14 * 8)

	// Cron2 represents clock daemon messages.
	Cron2 Facility = (15 * 8)

	// Local0 represents messages designated as local use 0.
	Local0 Facility = (16 * 8)

	// Local1 represents messages designated as local use 1.
	Local1 Facility = (17 * 8)

	// Local2 represents messages designated as local use 2.
	Local2 Facility = (18 * 8)

	// Local3 represents messages designated as local use 3.
	Local3 Facility = (19 * 8)

	// Local4 represents messages designated as local use 4.
	Local4 Facility = (20 * 8)

	// Local5 represents messages designated as local use 5.
	Local5 Facility = (21 * 8)

	// Local6 represents messages designated as local use 6.
	Local6 Facility = (22 * 8)

	// Local7 represents messages designated as local use 7.
	Local7 Facility = (23 * 8)

	InvalidFacility = errors.New(
		"A pri.Facility can only have one of the 24 values from" +
			" pri.Kern through pri.Local7. The acceptable values" +
			" are described at length by both RFC3164 and" +
			" RFC5424, and these constants are listed at:" +
			" https://godoc.org/github.com/proidiot/gone/log/pri",
	)
)

func (f Facility) Masked() Facility {
	// Equivalent to pri.Kern | pri.User | ... | pri.Local7
	mask := Facility(0xF8)
	return f & mask
}

func (f Facility) Valid() error {
	if f != f.Masked() || f > Local7 {
		return InvalidFacility
	} else {
		return nil
	}
}

func GetFromEnv() Facility {
	value := ""

	for _, env := range os.Environ() {
		// TODO add logic for checking if any facilities are vars
		switch env {
		case "LOG_FACILITY":
			value = os.Getenv(env)
		case "LOG_PRIORITY":
			if value == "" {
				value = os.Getenv(env)
			}
		}
	}

	switch value {
	case "LOG_KERN":
		return Kern
	case "LOG_USER":
		return User
	case "LOG_MAIL":
		return Mail
	case "LOG_DAEMON":
		return Daemon
	case "LOG_AUTH":
		return Auth
	case "LOG_SYSLOG":
		return Syslog
	case "LOG_LPR":
		return Lpr
	case "LOG_NEWS":
		return News
	case "LOG_UUCP":
		return Uucp
	case "LOG_CRON":
		return Cron
	case "LOG_AUTHPRIV":
		return Authpriv
	case "LOG_FTP":
		return Ftp
	case "LOG_NTP":
		return Ntp
	case "LOG_AUDIT":
		return Audit
	case "LOG_CONSOLE":
		return Console
	case "LOG_CRON2":
		return Cron2
	case "LOG_LOCAL0":
		return Local0
	case "LOG_LOCAL1":
		return Local1
	case "LOG_LOCAL2":
		return Local2
	case "LOG_LOCAL3":
		return Local3
	case "LOG_LOCAL4":
		return Local4
	case "LOG_LOCAL5":
		return Local5
	case "LOG_LOCAL6":
		return Local6
	case "LOG_LOCAL7":
		return Local7
	default:
		return User
	}
}

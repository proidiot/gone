//go:generate stringer -type Mask
package mask

import (
	"github.com/proidiot/gone/log/pri"
	"os"
)

// Log severity
// See RFC 3164 Section 4.1.1, RFC 5424 Section 6.2.1, and RFC 5427 Section 3.
type Mask byte

// Masks which only allow messages matching the corresponding severity level to
// be processed. These values would be passed to SetLogMask.
const (
	Emerg   Mask = (1 << pri.Emerg)
	Alert   Mask = (1 << pri.Alert)
	Crit    Mask = (1 << pri.Crit)
	Err     Mask = (1 << pri.Err)
	Warning Mask = (1 << pri.Warning)
	Notice  Mask = (1 << pri.Notice)
	Info    Mask = (1 << pri.Info)
	Debug   Mask = (1 << pri.Debug)
)

func UpTo(s pri.Severity) Mask {
	if s >= pri.Debug {
		return Mask(0xFF)
	} else {
		return Mask((1 << s) - 1)
	}
}

func (m Mask) Masked(s pri.Severity) bool {
	if s > pri.Debug {
		return true
	} else {
		return (m & (1 << s)) != 0
	}
}

func GetFromEnv() Mask {
	value := ""
	var m Mask
	upTo := false

	for _, env := range os.Environ() {
		// TODO read mask levels from env, union LOG_MASK
		switch env {
		case "LOG_UPTO":
			value = os.Getenv(env)
			upTo = true
		case "LOG_MASK":
			if value == "" {
				value = os.Getenv(env)
			}
		}
	}

	switch value {
	case "LOG_EMERG":
		m = Emerg
	case "LOG_ALERT":
		m = Alert
	case "LOG_CRIT":
		m = Crit
	case "LOG_ERR":
		m = Err
	case "LOG_WARNING":
		m = Warning
	case "LOG_NOTICE":
		m = Notice
	case "LOG_INFO":
		m = Info
	case "LOG_DEBUG":
		m = Debug
	default:
		m = Mask(0xFF)
	}

	if upTo {
		if m >= Debug {
			return Mask(0xFF)
		} else {
			return (m << 1) - 1
		}
	} else {
		return m
	}
}

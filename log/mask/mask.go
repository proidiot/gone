package mask

import (
	"fmt"
	"github.com/proidiot/gone/log/pri"
	"os"
	"strings"
)

// Log severity mask
// See POSIX.1-2008. Also see POSIX.1-2001, RFC 3164, and RFC 5424.
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

var lookup = map[Mask]string{
	Emerg:   "LOG_EMERG",
	Alert:   "LOG_ALERT",
	Crit:    "LOG_CRIT",
	Err:     "LOG_ERR",
	Warning: "LOG_WARNING",
	Notice:  "LOG_NOTICE",
	Info:    "LOG_INFO",
	Debug:   "LOG_DEBUG",
}

func UpTo(p pri.Priority) Mask {
	return Mask((1 << (p.Severity() + 1)) - 1)
}

func (m Mask) Masked(p pri.Priority) bool {
	return (m & (1 << p.Severity())) != 0
}

func (m Mask) String() string {
	if m == 0xFF {
		return "LOG_UPTO(LOG_DEBUG)"
	} else if m == 0x00 {
		return "LOG_MASK(0x0)"
	} else if _, present := lookup[m+1]; present {
		return fmt.Sprintf(
			"LOG_UPTO(%s)",
			lookup[(m+1)>>1],
		)
	} else {
		masked := []string{}
		for t := byte(1); t <= 0xFF && t != 0 && t <= byte(m); t *= 2 {
			if (t & byte(m)) != 0 {
				masked = append(masked, lookup[Mask(t)])
			}
		}
		return fmt.Sprintf("LOG_MASK(%s)", strings.Join(masked, "|"))
	}
}

func GetFromEnv() Mask {
	if val, set := os.LookupEnv("LOG_UPTO"); set {
		for mask, mstring := range lookup {
			if val == mstring {
				if mask == Debug {
					return Mask(0xFF)
				} else {
					return (mask << 1) - 1
				}
			}
		}

		// An invalid mask value was given, abort to default.
		return Mask(0xFF)
	} else if vals, set := os.LookupEnv("LOG_MASK"); set {
		m := Mask(0)

	maskLoop:
		for _, val := range strings.Split(vals, "|") {
			for mask, mstring := range lookup {
				if val == mstring {
					m |= mask
					continue maskLoop
				}
			}

			// An invalid mask value was given, abort to default.
			return Mask(0xFF)
		}

		return m
	} else {
		// No mask value given, use default.
		return Mask(0xFF)
	}
}

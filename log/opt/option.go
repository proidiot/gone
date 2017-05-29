package opt

import (
	"fmt"
	"os"
	"strings"
)

type Option byte

const (
	// Pid enables logging the process ID with each message. This is useful
	// for identifying specific processes.
	Pid Option = 0x01

	// Cons writes messages to the system console if they cannot be sent to
	// syslogd. This per-call fallback behavior is independent of the
	// stderr logging behavior described elsewhere.
	Cons Option = 0x02

	// ODelay delays openning resources until Syslog() is called. This
	// behavior is mutually exclusive with NDelay, but they are represented
	// by different options for historical reasons. This is the current
	// default behavior as mandated by POSIX, but it would be best to
	// explicitly set the option if there is a specific need for this
	// behavior as this particular implementation may one day stop using
	// this behavior by default (in fact, proidiot has already considered
	// making this change).
	ODelay Option = 0x04

	// NDelay opens the connection to the logging facility immediately.
	// Normally the open is delayed until the first message is logged. This
	// is useful for programs that need to manage the order in which file
	// descriptors are allocated. This behavior is mutually exclusive with
	// ODelay, but they are represented by different options for historical
	// reasons.
	NDelay Option = 0x08

	// NoWait enables the use of goroutines to send messages to syslogd so
	// that the calling function can proceed without waiting on these calls
	// to finish, although this hides any Syslog() errors.
	NoWait Option = 0x10

	// Perror prints all messages to stderr in addition to syslogd and the
	// system console (if Cons is set). This option isn't POSIX, but it is
	// very common. Since this option being set requires all messages to be
	// written to stderr (regardless of whether syslogd or the system
	// console also successfully received messages), this option being set
	// obviates the unusual fallback behavior of this particular
	// implementation, and so renders the NoFallback option useless.
	Perror Option = 0x20

	// Disable using stderr as a fallback when syslogd can't be reached.
	// This option isn't POSIX, and it is very unusual default behavior,
	// but proidiot likes it. This option is effectively meaningless if
	// Perror is set.
	NoFallback Option = 0x40
)

var lookup = map[Option]string{
	Pid:        "LOG_PID",
	Cons:       "LOG_CONS",
	ODelay:     "LOG_ODELAY",
	NDelay:     "LOG_NDELAY",
	NoWait:     "LOG_NOWAIT",
	Perror:     "LOG_PERROR",
	NoFallback: "LOG_NOFALLBACK",
}

func GetFromEnv() Option {
	var o Option

	for oflag, ostring := range lookup {
		if _, set := os.LookupEnv(ostring); set {
			o |= oflag
		}
	}

	return o
}

func (o Option) String() string {
	if o == 0 {
		return fmt.Sprintf("Option(%x)", 0)
	}

	oset := []string{}
	for t := byte(1); t <= 0x80 && t != 0x00 && t <= byte(o); t *= 2 {
		if (byte(o) & t) != 0 {
			if v, present := lookup[Option(t)]; present {
				oset = append(oset, v)
			} else {
				oset = append(
					oset,
					fmt.Sprintf(
						"Option(%x)",
						t,
					),
				)
			}
		}
	}

	return strings.Join(oset, "|")
}

package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/opt"
	"github.com/proidiot/gone/log/pri"
	"io"
	"os"
	"sync"
)

// Posixish is a syslogger.Syslogger that behaves much like the syslog system
// specified in POSIX.
type Posixish struct {
	i string
	o opt.Option
	f pri.Priority
	l Syslogger
	c []io.Closer
	x sync.RWMutex
}

var posixishNewNativeSyslog = NewNativeSyslog
var posixishOsOpen = os.Open
var posixishNewDelay = NewDelay
var posixishOsStderr = os.Stderr

// Syslog logs a message. How this message is routed depends on what settings
// were given to Openlog (and potentially a log mask).
func (px *Posixish) Syslog(p pri.Priority, msg interface{}) error {
	px.x.RLock()
	// The read unlock isn't being deferred here because t.Syslog could
	// take some time in the best case, but it might even need to acquire a
	// write lock if the syslog connection creation is being deferred.
	t := px.l
	px.x.RUnlock()

	if t == nil {
		px.x.Lock()

		px.f = pri.User

		e := px.prepareDelay()
		t = px.l

		px.x.Unlock()

		if e != nil {
			return e
		}
	}

	return t.Syslog(p, msg)
}

// Openlog re-initializes the Posixish based on the given opt.Option, overriding
// any previous values.
func (px *Posixish) Openlog(
	ident string,
	options opt.Option,
	facility pri.Priority,
) error {
	if (options&opt.NDelay) != 0 && (options&opt.ODelay) != 0 {
		return errors.New(
			"LOG_ODELAY and LOG_NDELAY are both being passed to" +
				" Openlog, but these options are mutually" +
				" exclusive.",
		)
	}

	if e := facility.ValidFacility(); e != nil {
		return e
	}

	px.x.Lock()
	defer px.x.Unlock()

	px.i = ident
	px.o = options
	px.f = facility

	if (px.o & opt.NDelay) != 0 {
		l, e := px.openlog()
		if e != nil {
			return e
		}

		px.l = l
		return nil
	}

	return px.prepareDelay()
}

// Close closes a Posixish, which has basically no effect other than to reset
// all the file descriptors.
func (px *Posixish) Close() error {
	return px.Closelog()
}

// Closelog closes a Posixish, which has basically no effect other than to reset
// all the file descriptors.
func (px *Posixish) Closelog() error {
	px.x.Lock()
	defer px.x.Unlock()
	return px.closelog()
}

// SetLogMask sets the Posixish's log mask.Mask.
func (px *Posixish) SetLogMask(m mask.Mask) error {
	px.x.Lock()
	defer px.x.Unlock()
	if px.l == nil {
		px.f = pri.User

		if e := px.prepareDelay(); e != nil {
			return e
		}
	}
	px.l = &SeverityMask{
		Syslogger: px.l,
		Mask:      m,
	}
	return nil
}

func (px *Posixish) prepareDelay() error {
	l, e := posixishNewDelay(
		func() (Syslogger, error) {
			px.x.Lock()
			defer px.x.Unlock()
			return px.openlog()
		},
	)

	if e != nil {
		return e
	}

	px.l = l
	return nil
}

func (px *Posixish) openlog() (Syslogger, error) {
	if e := px.closelog(); e != nil {
		return nil, e
	}

	var l Syslogger

	if n, e := posixishNewNativeSyslog(px.f, px.i); e == nil {
		px.c = append(px.c, n)
		l = n
	}

	if (px.o & opt.Cons) != 0 {
		if f, e := posixishOsOpen("/dev/console"); e == nil {
			px.c = append(px.c, f)

			c := &Rfc3164{
				Syslogger: &Writer{f},
				Facility:  px.f,
				Ident:     px.i,
				Pid:       (px.o & opt.Pid) != 0,
			}

			if l != nil {
				l = &Fallthrough{
					Default:     l,
					Fallthrough: c,
				}
			} else {
				l = c
			}
		}
	}

	if (px.o&opt.Perror) == 0 && (px.o&opt.NoFallback) != 0 {
		if l == nil {
			return nil, errors.New(
				"The posixish.Syslogger was unable to" +
					" connect to syslogd (and also" +
					" unable to connect to the system" +
					" console if that was requested)," +
					" but the NoFallback option has been" +
					" specified and the Perror option" +
					" has not been specified. As a" +
					" result, there is no mechanism for" +
					" recording logs.",
			)
		}
	} else {
		es := &Rfc3164{
			Syslogger: &Writer{posixishOsStderr},
			Facility:  px.f,
			Ident:     px.i,
			Pid:       (px.o & opt.Pid) != 0,
		}

		if l == nil {
			l = es
		} else if (px.o & opt.Perror) != 0 {
			l = &Multi{
				Sysloggers: []Syslogger{
					l,
					es,
				},
				TryAll: true,
			}
		} else {
			l = &Fallthrough{
				Default:     l,
				Fallthrough: es,
			}
		}
	}

	if (px.o & opt.NoWait) != 0 {
		l = &NoWait{l}
	}

	return l, nil
}

func (px *Posixish) closelog() error {
	var err error

	for _, c := range px.c {
		e := c.Close()
		if err == nil && e != nil {
			err = e
		}
	}

	px.l = nil

	return err
}

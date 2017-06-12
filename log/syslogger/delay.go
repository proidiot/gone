package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"sync"
)

type sysloggerHandle struct {
	s Syslogger
}

// Delay is a syslogger.Syslogger that delays the initialization of another
// syslogger.Syslogger until the first time that Syslog is called.
type Delay struct {
	h  *sysloggerHandle
	cb func() (Syslogger, error)
	x  sync.Mutex
}

// Syslog logs messages. In the case of Delay, these messages are passed
// unaltered to another syslogger.Syslogger, and that syslogger.Syslogger would
// be created at this point if it had not already existed.
func (d *Delay) Syslog(p pri.Priority, msg interface{}) error {
	d.x.Lock()
	// Not deferring the unlock here because the actual Syslog call at the
	// end may take some time.
	if d.h == nil {
		s, e := d.cb()
		if e != nil {
			d.x.Unlock()
			return e
		}

		d.h = &sysloggerHandle{s}
	}
	h := d.h
	d.x.Unlock()
	return h.s.Syslog(p, msg)
}

// NewDelay gives a Delay syslogger.Syslogger given the callback function which
// will ultimately be used to create the real syslogger.Syslogger to be used.
func NewDelay(cb func() (Syslogger, error)) (*Delay, error) {
	if cb == nil {
		return nil, errors.New(
			"A syslogger.Delay must have a valid callback for" +
				" generating a Syslogger, but a nil callback" +
				" was given to syslogger.NewDelay(...).",
		)
	}
	return &Delay{cb: cb}, nil
}

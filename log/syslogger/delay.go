package syslogger

import (
	"github.com/proidiot/gone/log/pri"
	"sync"
)

type sysloggerHandle struct {
	s Syslogger
}

type Delay struct {
	h *sysloggerHandle
	cb func() (Syslogger, error)
	x *sync.Mutex
}

func (d Delay) Syslog(p pri.Priority, msg interface{}) error {
	d.x.Lock()
	// Not deferring the unlock here because the actual Syslog call at the
	// end may take some time.
	if d.h == nil {
		if s, e := d.cb(); e != nil {
			d.x.Unlock()
			return e
		} else {
			d.h = &sysloggerHandle{s}
		}
	}
	h := d.h
	d.x.Unlock()
	return h.s.Syslog(p, msg)
}

func (d Delay) Reset() {
	d.x.Lock()
	defer d.x.Unlock()
	d.h = nil
}

func NewDelay(cb func() (Syslogger, error)) Delay {
	return Delay{
		cb: cb,
	}
}

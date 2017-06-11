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

func (p *Posixish) Openlog(
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

	p.x.Lock()
	defer p.x.Unlock()

	p.i = ident
	p.o = options
	p.f = facility

	if (p.o & opt.NDelay) != 0 {
		if l, e := p.openlog(); e != nil {
			return e
		} else {
			p.l = l
			return nil
		}
	} else {
		return p.prepareDelay()
	}
}

func (p *Posixish) Close() error {
	return p.Closelog()
}

func (p *Posixish) Closelog() error {
	p.x.Lock()
	defer p.x.Unlock()
	return p.closelog()
}

func (p *Posixish) SetLogMask(m mask.Mask) error {
	p.x.Lock()
	defer p.x.Unlock()
	if p.l == nil {
		p.f = pri.User

		if e := p.prepareDelay(); e != nil {
			return e
		}
	}
	p.l = &SeverityMask{
		Syslogger: p.l,
		Mask:      m,
	}
	return nil
}

func (p *Posixish) prepareDelay() error {
	l, e := posixishNewDelay(
		func() (Syslogger, error) {
			p.x.Lock()
			defer p.x.Unlock()
			return p.openlog()
		},
	)

	if e != nil {
		return e
	} else {
		p.l = l

		return nil
	}
}

func (p *Posixish) openlog() (Syslogger, error) {
	if e := p.closelog(); e != nil {
		return nil, e
	}

	var l Syslogger

	if n, e := posixishNewNativeSyslog(p.f, p.i); e == nil {
		p.c = append(p.c, n)
		l = n
	}

	if (p.o & opt.Cons) != 0 {
		if f, e := posixishOsOpen("/dev/console"); e == nil {
			p.c = append(p.c, f)

			c := &Rfc3164{
				Syslogger: &Writer{f},
				Facility:  p.f,
				Ident:     p.i,
				Pid:       (p.o & opt.Pid) != 0,
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

	if (p.o&opt.Perror) == 0 && (p.o&opt.NoFallback) != 0 {
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
			Facility:  p.f,
			Ident:     p.i,
			Pid:       (p.o & opt.Pid) != 0,
		}

		if l == nil {
			l = es
		} else if (p.o & opt.Perror) != 0 {
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

	if (p.o & opt.NoWait) != 0 {
		l = &NoWait{l}
	}

	return l, nil
}

func (p *Posixish) closelog() error {
	var err error

	for _, c := range p.c {
		e := c.Close()
		if err == nil && e != nil {
			err = e
		}
	}

	p.l = nil

	return err
}

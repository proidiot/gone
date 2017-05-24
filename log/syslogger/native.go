package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"log/syslog"
)

type NativeSyslog struct {
	w *syslog.Writer
	f pri.Facility
}

func (n NativeSyslog) Syslog(p pri.Priority, msg interface{}) error {
	m, goodType := msg.(string)
	if !goodType {
		return errors.New(
			"The native Go log/syslog system only accepts" +
			" strings as a message, but a non-string message was" +
			" given.",
		)
	}

	if p.Facility != 0x00 && p.Facility != n.f {
		return errors.New(
			fmt.Sprintf(
				"The native Go log/syslog system does not" +
				" provide a mechanism for changing log" +
				" facilities of an existing *syslog.Writer," +
				" but the pri.Facility this" +
				" *syslogger.NativeSyslog was created with" +
				" does not match the pri.Facility component" +
				" of the given pr.Priority. This" +
				" *syslogger.NativeSyslog was created with" +
				" pri.Facility %s, but the given" +
				" pri.Priority argument has pri.Facility: %s",
				n.f,
				p.Facility,
			),
		)
	}

	switch p.Severity {
	case pri.Emerg:
		return n.w.Emerg(m)
	case pri.Alert:
		return n.w.Alert(m)
	case pri.Crit:
		return n.w.Crit(m)
	case pri.Err:
		return n.w.Err(m)
	case pri.Warning:
		return n.w.Warning(m)
	case pri.Notice:
		return n.w.Notice(m)
	case pri.Info:
		return n.w.Info(m)
	case pri.Debug:
		return n.w.Debug(m)
	default:
		return errors.New(
			fmt.Sprintf(
				"The given pri.Priority argument has an" +
				" invalid pri.Severity component: %s",
				p.Severity,
			),
		)
	}
}

func NewNativeSyslog(p pri.Priority, ident string) (NativeSyslog, error) {
	if e := p.Facility.Valid(); e != nil {
		return nil, e
	}

	w, e := syslog.New(syslog.Priority(p.Facility.Masked()), ident)
	if e != nil {
		return nil, e
	}

	return NativeSyslog{
		w,
		p.Facility,
	}, nil
}

func DialNativeSyslog(
	network string,
	raddr string,
	p pri.Priority,
	ident string,
) (NativeSyslog, error) {
	if e := p.Facility.Valid(); e != nil {
		return nil, e
	}

	w, e := syslog.Dial(
		network,
		raddr,
		syslog.Priority(p.Facility.Masked()),
		ident,
	)
	if e != nil {
		return nil, e
	}

	return NativeSyslog{
		w,
		p.Facility,
	}, nil
}

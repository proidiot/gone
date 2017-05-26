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
				" strings as a message, but a non-string" +
				" message was given.",
		)
	}

	if p.Facility != 0x00 && p.Facility != n.f {
		return errors.New(
			fmt.Sprintf(
				"The native Go log/syslog system does not"+
					" provide a mechanism for changing"+
					" log facilities of an existing"+
					" *syslog.Writer, but the"+
					" pri.Facility this"+
					" *syslogger.NativeSyslog was"+
					" created with does not match the"+
					" pri.Facility component of the"+
					" given pr.Priority. This"+
					" *syslogger.NativeSyslog was"+
					" created with pri.Facility %s, but"+
					" the given pri.Priority argument"+
					" has pri.Facility: %s",
				n.f,
				p.Facility,
			),
		)
	}

	switch p.Severity.Masked() {
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
				"The given pri.Priority argument has an"+
					" invalid pri.Severity component: %s",
				p.Severity,
			),
		)
	}
}

func (n NativeSyslog) Close() error {
	return n.w.Close()
}

func NewNativeSyslog(f pri.Facility, ident string) (NativeSyslog, error) {
	if e := f.Valid(); e != nil {
		return NativeSyslog{}, e
	}

	w, e := syslog.New(syslog.Priority(f.Masked()), ident)
	if e != nil {
		return NativeSyslog{}, e
	}

	return NativeSyslog{
		w,
		f,
	}, nil
}

func DialNativeSyslog(
	network string,
	raddr string,
	f pri.Facility,
	ident string,
) (NativeSyslog, error) {
	if e := f.Valid(); e != nil {
		return NativeSyslog{}, e
	}

	w, e := syslog.Dial(
		network,
		raddr,
		syslog.Priority(f.Masked()),
		ident,
	)
	if e != nil {
		return NativeSyslog{}, e
	}

	return NativeSyslog{
		w,
		f,
	}, nil
}

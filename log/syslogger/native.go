package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"log/syslog"
)

// NativeSyslog is a syslogger.Syslogger that is a lightweight wrapper around
// golang's native log/syslog.Writer.
type NativeSyslog struct {
	w *syslog.Writer
	f pri.Priority
}

// Syslog logs a message. In the case of NativeSyslog, the message is sent to
// golang's log/syslog.Writer.
func (n *NativeSyslog) Syslog(p pri.Priority, msg interface{}) error {
	m, goodType := msg.(string)
	if !goodType {
		return errors.New(
			"The native Go log/syslog system only accepts" +
				" strings as a message, but a non-string" +
				" message was given.",
		)
	}

	if p.Facility() != 0x00 && p.Facility() != n.f {
		return fmt.Errorf(
			"The native Go log/syslog system does not provide a"+
				" mechanism for changing log facilities of an"+
				" existing *syslog.Writer, but the"+
				" pri.Facility this *syslogger.NativeSyslog"+
				" was created with does not match the"+
				" pri.Facility component of the given"+
				" pr.Priority. This *syslogger.NativeSyslog"+
				" was created with pri.Facility %s, but the"+
				" given pri.Priority argument has"+
				" pri.Facility: %s",
			n.f,
			p.Facility(),
		)
	}

	switch p.Severity() {
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
	default:
		return n.w.Debug(m)
	}
}

// Close closes a native log/syslog.Writer.
func (n *NativeSyslog) Close() error {
	return n.w.Close()
}

// NewNativeSyslog creates a new NativeSyslog based on the given log facility
// and identity string.
func NewNativeSyslog(f pri.Priority, ident string) (*NativeSyslog, error) {
	if e := f.ValidFacility(); e != nil {
		return nil, e
	}

	w, e := syslogNew(syslog.Priority(f.Facility()), ident)
	if e != nil {
		return nil, e
	}

	return &NativeSyslog{
		w,
		f,
	}, nil
}

// DialNativeSyslog creates a new NativeSyslog based on using log/syslog.Dial.
func DialNativeSyslog(
	network string,
	raddr string,
	f pri.Priority,
	ident string,
) (*NativeSyslog, error) {
	if e := f.ValidFacility(); e != nil {
		return nil, e
	}

	w, e := syslog.Dial(
		network,
		raddr,
		syslog.Priority(f.Facility()),
		ident,
	)
	if e != nil {
		return nil, e
	}

	return &NativeSyslog{
		w,
		f,
	}, nil
}

package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/opt"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/syslog"
	"os"
	"testing"
)

func TestPosixishOpenlog(t *testing.T) {
	origNewNativeSyslog := posixishNewNativeSyslog
	defer func() {
		posixishNewNativeSyslog = origNewNativeSyslog
	}()
	fakeNewNativeSyslog := func(
		f pri.Priority,
		s string,
	) (*NativeSyslog, error) {
		return &NativeSyslog{
			w: &syslog.Writer{},
			f: f,
		}, nil
	}
	errorNewNativeSyslog := func(
		pri.Priority,
		string,
	) (*NativeSyslog, error) {
		return nil, errors.New("Artificial error for NewNativeSyslog")
	}

	origOsOpen := posixishOsOpen
	defer func() {
		posixishOsOpen = origOsOpen
	}()
	fakeOsOpen := func(string) (*os.File, error) {
		return &os.File{}, nil
	}
	errorOsOpen := func(string) (*os.File, error) {
		return nil, errors.New("Artificial error for os.Open")
	}

	origNewDelay := posixishNewDelay
	defer func() {
		posixishNewDelay = origNewDelay
	}()
	errorNewDelay := func(func() (Syslogger, error)) (*Delay, error) {
		return nil, errors.New("Artificial error for NewDelay")
	}

	devNull, e := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	require.NoError(
		t,
		e,
		"Posixish Syslog test requires the ability to open /dev/null"+
			" in append mode.",
	)
	origOsStderr := posixishOsStderr
	defer func() {
		_ = devNull.Close()
		posixishOsStderr = origOsStderr
	}()
	posixishOsStderr = devNull

	type testCase struct {
		inputIdent                string
		inputOptions              opt.Option
		inputFacility             pri.Priority
		causeNewNativeSyslogError bool
		causeOsOpenError          bool
		causeNewDelayError        bool
		expectedError             bool
		expectedSysloggerType     Syslogger
		expectedClosers           []io.Closer
	}

	tests := map[string]testCase{
		"nil values": {
			inputIdent:            "",
			inputOptions:          opt.Option(0x0),
			inputFacility:         pri.Priority(0x0),
			expectedError:         false,
			expectedSysloggerType: &Delay{},
		},
		"priority collision": {
			inputOptions:  opt.NDelay | opt.ODelay,
			expectedError: true,
		},
		"bad facility": {
			inputFacility: pri.Priority(0xF8),
			expectedError: true,
		},
		"severity as facility": {
			inputFacility: pri.Warning,
			expectedError: true,
		},
		"no delay": {
			inputOptions:          opt.NDelay,
			expectedError:         false,
			expectedSysloggerType: &Fallthrough{},
			expectedClosers: []io.Closer{
				&NativeSyslog{},
			},
		},
		"no delay, native error": {
			inputOptions:              opt.NDelay,
			causeNewNativeSyslogError: true,
			expectedError:             false,
			expectedSysloggerType:     &Rfc3164{},
			expectedClosers:           []io.Closer{},
		},
		"no delay, cons": {
			inputOptions:          opt.NDelay | opt.Cons,
			expectedError:         false,
			expectedSysloggerType: &Fallthrough{},
			expectedClosers: []io.Closer{
				&NativeSyslog{},
				&os.File{},
			},
		},
		"no delay, cons, native error": {
			inputOptions:              opt.NDelay | opt.Cons,
			causeNewNativeSyslogError: true,
			expectedError:             false,
			expectedSysloggerType:     &Fallthrough{},
			expectedClosers: []io.Closer{
				&NativeSyslog{},
			},
		},
		"no delay, cons, cons error": {
			inputOptions:          opt.NDelay | opt.Cons,
			causeOsOpenError:      true,
			expectedError:         false,
			expectedSysloggerType: &Fallthrough{},
			expectedClosers: []io.Closer{
				&NativeSyslog{},
			},
		},
		"no delay, no fallback": {
			inputOptions:          opt.NDelay | opt.NoFallback,
			expectedError:         false,
			expectedSysloggerType: &NativeSyslog{},
			expectedClosers: []io.Closer{
				&NativeSyslog{},
			},
		},
		"no delay, no fallback, native error": {
			inputOptions:              opt.NDelay | opt.NoFallback,
			causeNewNativeSyslogError: true,
			expectedError:             true,
		},
		"no delay, perror": {
			inputOptions:          opt.NDelay | opt.Perror,
			expectedError:         false,
			expectedSysloggerType: &Multi{},
			expectedClosers: []io.Closer{
				&NativeSyslog{},
			},
		},
		"no delay, no wait": {
			inputOptions:          opt.NDelay | opt.NoWait,
			expectedError:         false,
			expectedSysloggerType: &NoWait{},
			expectedClosers: []io.Closer{
				&NativeSyslog{},
			},
		},
		"delay error": {
			causeNewDelayError:    true,
			expectedError:         true,
			expectedSysloggerType: nil,
		},
	}

	for explanation, test := range tests {
		p := new(Posixish)

		if test.causeNewNativeSyslogError {
			posixishNewNativeSyslog = errorNewNativeSyslog
		} else {
			posixishNewNativeSyslog = fakeNewNativeSyslog
		}

		if test.causeOsOpenError {
			posixishOsOpen = errorOsOpen
		} else {
			posixishOsOpen = fakeOsOpen
		}

		if test.causeNewDelayError {
			posixishNewDelay = errorNewDelay
		} else {
			posixishNewDelay = origNewDelay
		}

		actualError := p.Openlog(
			test.inputIdent,
			test.inputOptions,
			test.inputFacility,
		)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"Posixish Openlog test expects an error for:"+
					" %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"Posixish Openlog test expects no error for:"+
					" %s",
				explanation,
			)
		}

		actualSyslogger := p.l
		assert.IsType(
			t,
			test.expectedSysloggerType,
			actualSyslogger,
			"Posixish Openlog test expects a specific syslogger"+
				" for: %s",
			explanation,
		)

		actualClosers := p.c
		// TODO better closers test
		// odd that they swapped the order for length...
		assert.Len(
			t,
			actualClosers,
			len(test.expectedClosers),
			"Posixish Openlog test expects specific closers for:"+
				" %s",
			explanation,
		)

		// Closing an empty file will cause errors we can ignore.
		_ = p.Close()
	}
}

func TestPosixishSyslog(t *testing.T) {
	origNewNativeSyslog := posixishNewNativeSyslog
	defer func() {
		posixishNewNativeSyslog = origNewNativeSyslog
	}()
	errorNewNativeSyslog := func(
		pri.Priority,
		string,
	) (*NativeSyslog, error) {
		return nil, errors.New("Artificial error for NewNativeSyslog")
	}
	posixishNewNativeSyslog = errorNewNativeSyslog

	origNewDelay := posixishNewDelay
	defer func() {
		posixishNewDelay = origNewDelay
	}()
	errorNewDelay := func(func() (Syslogger, error)) (*Delay, error) {
		return nil, errors.New("Artificial error for NewDelay")
	}

	devNull, e := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	require.NoError(
		t,
		e,
		"Posixish Syslog test requires the ability to open /dev/null"+
			" in append mode.",
	)
	origOsStderr := posixishOsStderr
	defer func() {
		_ = devNull.Close()
		posixishOsStderr = origOsStderr
	}()
	posixishOsStderr = devNull

	type testCase struct {
		inputPriority      pri.Priority
		inputMsg           interface{}
		causeOpenlogError  bool
		causeBlankOpenlog  bool
		causeNewDelayError bool
		expectedError      bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputPriority: pri.Priority(0x0),
			inputMsg:      nil,
			expectedError: true,
		},
		"full values": {
			inputPriority: pri.Priority(0x0),
			inputMsg:      "full values",
			expectedError: false,
		},
		"openlog error": {
			inputPriority:     pri.Info,
			inputMsg:          "openlog error msg",
			causeOpenlogError: true,
			expectedError:     true,
		},
		"err": {
			inputPriority: pri.Err,
			inputMsg:      "err msg",
			expectedError: false,
		},
		"blank Openlog first": {
			inputPriority:     pri.Crit,
			inputMsg:          "blank Openlog first msg",
			causeBlankOpenlog: true,
			expectedError:     false,
		},
		"delay error": {
			inputPriority:      pri.Notice,
			inputMsg:           "delay error msg",
			causeNewDelayError: true,
			expectedError:      true,
		},
	}

	for explanation, test := range tests {
		p := new(Posixish)

		if test.causeBlankOpenlog {
			openError := p.Openlog(
				"",
				opt.Option(0x0),
				pri.Priority(0x0),
			)
			assert.NoError(
				t,
				openError,
				"Posixish Syslog test expects no error during"+
					" Openlog for: %s",
				explanation,
			)
		}
		if test.causeOpenlogError {
			p.o |= opt.NoFallback
		}
		if test.causeNewDelayError {
			posixishNewDelay = errorNewDelay
		} else {
			posixishNewDelay = origNewDelay
		}

		actualError := p.Syslog(test.inputPriority, test.inputMsg)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"Posixish Syslog test expects an error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"Posixish Syslog test expects no error for: %s",
				explanation,
			)
		}
	}
}

func TestPosixishSetLogMask(t *testing.T) {
	origNewDelay := posixishNewDelay
	defer func() {
		posixishNewDelay = origNewDelay
	}()
	errorNewDelay := func(func() (Syslogger, error)) (*Delay, error) {
		return nil, errors.New("Artificial error for NewDelay")
	}

	type testCase struct {
		inputMask               mask.Mask
		causeBlankOpenlog       bool
		causeNewDelayError      bool
		expectedError           bool
		expectedMaskedSyslogger Syslogger
	}

	tests := map[string]testCase{
		"nil values": {
			inputMask:               mask.Mask(0x0),
			expectedError:           false,
			expectedMaskedSyslogger: &Delay{},
		},
		"follow-on log mask": {
			inputMask:               mask.Warning,
			causeBlankOpenlog:       true,
			expectedError:           false,
			expectedMaskedSyslogger: &Delay{},
		},
		"delay error": {
			inputMask:          mask.Notice,
			causeNewDelayError: true,
			expectedError:      true,
		},
	}

	for explanation, test := range tests {
		p := new(Posixish)

		if test.causeBlankOpenlog {
			openError := p.Openlog(
				"",
				opt.Option(0x0),
				pri.Priority(0x0),
			)
			assert.NoError(
				t,
				openError,
				"Posixish SetLogMask test expects no error"+
					" during Openlog for: %s",
				explanation,
			)
		}
		if test.causeNewDelayError {
			posixishNewDelay = errorNewDelay
		} else {
			posixishNewDelay = origNewDelay
		}

		actualError := p.SetLogMask(test.inputMask)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"Posixish SetLogMask test expects an error"+
					" for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"Posixish SetLogMask test expects no error"+
					" for: %s",
				explanation,
			)
		}

		if actualError == nil {
			actualSyslogger := p.l
			assert.IsType(
				t,
				&SeverityMask{},
				actualSyslogger,
				"Posixish SetLogMask test expects a"+
					" *syslogger.SeverityMask as the"+
					" Posixish syslogger for: %s",
				explanation,
			)

			if s, ok := actualSyslogger.(*SeverityMask); ok {
				actualMaskedSyslogger := s.Syslogger
				assert.IsType(
					t,
					test.expectedMaskedSyslogger,
					actualMaskedSyslogger,
					"Posixish SetLogMask expects a"+
						" different masked syslogger"+
						" for: %s",
					explanation,
				)
			}
		}
	}
}

func TestPosixishCloseError(t *testing.T) {
	p := new(Posixish)

	p.c = append(p.c, &errorCloser{})

	actualError := p.Openlog("", opt.NDelay, pri.Priority(0x0))

	assert.Errorf(
		t,
		actualError,
		"Posixish Close error test expects an error.",
	)
}

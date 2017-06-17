package log

import (
	"testing"

	"github.com/proidiot/gone/log/opt"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
)

func TestOpenlog(t *testing.T) {
	type testCase struct {
		inputIdent          string
		inputOptions        opt.Option
		inputFacility       pri.Priority
		useLimitedSyslogger bool
		expectedError       bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputIdent:    "",
			inputOptions:  opt.Option(0x0),
			inputFacility: pri.Priority(0x0),
			expectedError: false,
		},
		"limited syslogger": {
			inputIdent:          "limited",
			useLimitedSyslogger: true,
			expectedError:       true,
		},
	}

	for explanation, test := range tests {
		if test.useLimitedSyslogger {
			SetSyslogger(new(limitedSyslogger))
		} else {
			SetSyslogger(new(testSyslogger))
		}

		actualError := Openlog(
			test.inputIdent,
			test.inputOptions,
			test.inputFacility,
		)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"Openlog test expects an error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"Openlog test expects no error for: %s",
				explanation,
			)
		}

		closeError := Closelog()
		if !test.useLimitedSyslogger {
			assert.NoError(
				t,
				closeError,
				"Openlog test expects no error on Closelog"+
					" for: %s",
				explanation,
			)
		}
	}
}

func TestSyslogWrappers(t *testing.T) {
	type testCase struct {
		callFunc         func(interface{}) error
		causeSyslogError bool
		expectedError    bool
		expectedPriority pri.Priority
	}

	tests := map[string]testCase{
		"emerg": {
			callFunc:         Emerg,
			expectedError:    false,
			expectedPriority: pri.Emerg,
		},
		"emergency": {
			callFunc:         Emergency,
			expectedError:    false,
			expectedPriority: pri.Emerg,
		},
		"alert": {
			callFunc:         Alert,
			expectedError:    false,
			expectedPriority: pri.Alert,
		},
		"crit": {
			callFunc:         Crit,
			expectedError:    false,
			expectedPriority: pri.Crit,
		},
		"critical": {
			callFunc:         Critical,
			expectedError:    false,
			expectedPriority: pri.Crit,
		},
		"err": {
			callFunc:         Err,
			expectedError:    false,
			expectedPriority: pri.Err,
		},
		"error": {
			callFunc:         Error,
			expectedError:    false,
			expectedPriority: pri.Err,
		},
		"warning": {
			callFunc:         Warning,
			expectedError:    false,
			expectedPriority: pri.Warning,
		},
		"warn": {
			callFunc:         Warn,
			expectedError:    false,
			expectedPriority: pri.Warning,
		},
		"notice": {
			callFunc:         Notice,
			expectedError:    false,
			expectedPriority: pri.Notice,
		},
		"info": {
			callFunc:         Info,
			expectedError:    false,
			expectedPriority: pri.Info,
		},
		"information": {
			callFunc:         Information,
			expectedError:    false,
			expectedPriority: pri.Info,
		},
		"debug": {
			callFunc:         Debug,
			expectedError:    false,
			expectedPriority: pri.Debug,
		},
	}

	for explanation, test := range tests {
		s := new(testSyslogger)

		if test.causeSyslogError {
			s.TriggerError = true
		}

		SetSyslogger(s)

		expectedMsg := explanation

		actualError := test.callFunc(expectedMsg)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"Syslog wrapper test expects an error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"Syslog wrapper test expects no error for: %s",
				explanation,
			)
		}

		actualPriority := s.LastPri

		assert.Equal(
			t,
			test.expectedPriority,
			actualPriority,
			"Syslog wrapper test expects the priority to match"+
				" for: %s",
			explanation,
		)

		actualMsg := s.LastMsg

		assert.Equal(
			t,
			expectedMsg,
			actualMsg,
			"Syslog wrapper test expects the message to match for:"+
				" %s",
			explanation,
		)

		closeError := Closelog()
		if !test.causeSyslogError {
			assert.NoError(
				t,
				closeError,
				"Syslog wrapper test expects Closelog to have"+
					" no error for: %s",
				explanation,
			)
		}
	}
}

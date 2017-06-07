package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"log/syslog"
	"testing"
)

func TestNewNativeSyslog(t *testing.T) {
	var origSyslogNew = syslogNew
	defer func() {
		syslogNew = origSyslogNew
	}()

	tryRealSyslog := true
	if _, e := syslogNew(0, ""); e != nil {
		t.Log(
			"It seems that log/syslog.New is failing by default." +
				" This could be because the system doesn't" +
				" have a running syslogd (such as in a" +
				" container), or it could be because the" +
				" system isn't Unix-like. In order to" +
				" complete the tests, a mock of" +
				" log/syslog.New will be used for the happy" +
				" path in addition to the mock already" +
				" being used for the sad path.",
		)
		tryRealSyslog = false
	}

	type testCase struct {
		inputPriority       pri.Priority
		inputIdent          string
		causeSyslogNewError bool
		expectedError       bool
		expectedSyslogger   bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputPriority:     pri.Priority(0x0),
			inputIdent:        "",
			expectedError:     false,
			expectedSyslogger: true,
		},
		"full values": {
			inputPriority: pri.Priority(0xFF),
			inputIdent:    "full values ident",
			expectedError: true,
		},
		"broken syslog.New": {
			causeSyslogNewError: true,
			expectedError:       true,
		},
	}

	for explanation, test := range tests {
		if test.causeSyslogNewError {
			syslogNew = func(
				syslog.Priority,
				string,
			) (*syslog.Writer, error) {
				return nil, errors.New(
					"Artificial error for syslog.New",
				)
			}
		} else if tryRealSyslog {
			syslogNew = origSyslogNew
		} else {
			syslogNew = func(
				syslog.Priority,
				string,
			) (*syslog.Writer, error) {
				return &syslog.Writer{}, nil
			}
		}

		actualSyslogger, actualError := NewNativeSyslog(
			test.inputPriority,
			test.inputIdent,
		)

		if test.expectedError {
			assert.Error(
				t,
				actualError,
				fmt.Sprintf(
					"NewNativeSyslog test expects an"+
						" error for: %s",
					explanation,
				),
			)
		} else {
			assert.NoError(
				t,
				actualError,
				fmt.Sprintf(
					"NewNativeSyslog test expects no"+
						" error for: %s",
					explanation,
				),
			)
		}

		if test.expectedSyslogger {
			assert.NotNil(
				t,
				actualSyslogger,
				fmt.Sprintf(
					"NewNativeSyslog test expects non-nil"+
						" syslogger for: %s",
					explanation,
				),
			)
		} else {
			assert.Nil(
				t,
				actualSyslogger,
				fmt.Sprintf(
					"NewNativeSyslog test expects nil"+
						" syslogger for: %s",
					explanation,
				),
			)
		}
	}
}

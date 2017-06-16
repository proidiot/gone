package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoWaitSyslog(t *testing.T) {
	sfs, sc := newSyncFlagSyslogger()
	defer close(sc)

	type testCase struct {
		inputSyslogger Syslogger
		expectedError  bool
		expectedCall   bool
	}

	tests := map[string]testCase{
		"nil value": {
			inputSyslogger: nil,
			expectedError:  true,
		},
		"error syslogger": {
			inputSyslogger: &errorSyslogger{},
			expectedError:  false,
		},
		"flag syslogger": {
			inputSyslogger: sfs,
			expectedError:  false,
			expectedCall:   true,
		},
	}

	for explanation, test := range tests {
		n := &NoWait{
			Syslogger: test.inputSyslogger,
		}

		actualError := n.Syslog(pri.Priority(0x0), nil)

		if test.expectedError {
			assert.Error(
				t,
				actualError,
				fmt.Sprintf(
					"NoWait test expected error for: %s",
					explanation,
				),
			)
		} else {
			assert.NoError(
				t,
				actualError,
				fmt.Sprintf(
					"NoWait test unexpected error for: %s",
					explanation,
				),
			)
		}

		if test.inputSyslogger == sfs {
			sc <- nil
			actualCall := sfs.Flag

			assert.Equal(
				t,
				test.expectedCall,
				actualCall,
				fmt.Sprintf(
					"NoWait test call check failure for:"+
						" %s",
					explanation,
				),
			)
		}
	}
}

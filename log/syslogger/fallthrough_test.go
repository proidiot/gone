package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFallthroughSyslog(t *testing.T) {
	type testCase struct {
		inputDefault            Syslogger
		inputFallthrough        Syslogger
		expectedError           bool
		expectedDefaultCall     bool
		expectedFallthroughCall bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputDefault:     nil,
			inputFallthrough: nil,
			expectedError:    true,
		},
		"error default, nil fallthrough": {
			inputDefault:     &errorSyslogger{},
			inputFallthrough: nil,
			expectedError:    true,
		},
		"working default, nil fallthrough": {
			inputDefault:        &flagSyslogger{},
			inputFallthrough:    nil,
			expectedError:       false,
			expectedDefaultCall: true,
		},
		"nil default, error fallthrough": {
			inputDefault:     nil,
			inputFallthrough: &errorSyslogger{},
			expectedError:    true,
		},
		"error default, error fallthrough": {
			inputDefault:     &errorSyslogger{},
			inputFallthrough: &errorSyslogger{},
			expectedError:    true,
		},
		"working default, error fallthrough": {
			inputDefault:        &flagSyslogger{},
			inputFallthrough:    &errorSyslogger{},
			expectedError:       false,
			expectedDefaultCall: true,
		},
		"nil default, working fallthrough": {
			inputDefault:            nil,
			inputFallthrough:        &flagSyslogger{},
			expectedError:           false,
			expectedFallthroughCall: true,
		},
		"error default, working fallthrough": {
			inputDefault:            &errorSyslogger{},
			inputFallthrough:        &flagSyslogger{},
			expectedError:           false,
			expectedFallthroughCall: true,
		},
		"working default, working fallthrough": {
			inputDefault:            &flagSyslogger{},
			inputFallthrough:        &flagSyslogger{},
			expectedError:           false,
			expectedDefaultCall:     true,
			expectedFallthroughCall: false,
		},
	}

	for explanation, test := range tests {
		f := &Fallthrough{
			Default:     test.inputDefault,
			Fallthrough: test.inputFallthrough,
		}

		actualError := f.Syslog(pri.Priority(0x0), nil)

		if test.expectedError {
			assert.Error(
				t,
				actualError,
				fmt.Sprintf(
					"Fallthrough test expected error"+
						" for: %s",
					explanation,
				),
			)
		} else {
			assert.NoError(
				t,
				actualError,
				fmt.Sprintf(
					"Fallthrough test unexpected error"+
						" for: %s",
					explanation,
				),
			)
		}

		if fs, ok := test.inputDefault.(*flagSyslogger); ok {
			actualCall := fs.Flag

			assert.Equal(
				t,
				test.expectedDefaultCall,
				actualCall,
				fmt.Sprintf(
					"Fallthrough test call check failure"+
						" on default syslogger for: %s",
					explanation,
				),
			)
		}

		if fs, ok := test.inputFallthrough.(*flagSyslogger); ok {
			actualCall := fs.Flag

			assert.Equal(
				t,
				test.expectedFallthroughCall,
				actualCall,
				fmt.Sprintf(
					"Fallthrough test call check failure"+
						" on fallthrough syslogger"+
						" for: %s",
					explanation,
				),
			)
		}
	}
}

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
			inputDefault:     &errorSyslog{},
			inputFallthrough: nil,
			expectedError:    true,
		},
		"working default, nil fallthrough": {
			inputDefault:        &flagSyslog{},
			inputFallthrough:    nil,
			expectedError:       false,
			expectedDefaultCall: true,
		},
		"nil default, error fallthrough": {
			inputDefault:     nil,
			inputFallthrough: &errorSyslog{},
			expectedError:    true,
		},
		"error default, error fallthrough": {
			inputDefault:     &errorSyslog{},
			inputFallthrough: &errorSyslog{},
			expectedError:    true,
		},
		"working default, error fallthrough": {
			inputDefault:        &flagSyslog{},
			inputFallthrough:    &errorSyslog{},
			expectedError:       false,
			expectedDefaultCall: true,
		},
		"nil default, working fallthrough": {
			inputDefault:            nil,
			inputFallthrough:        &flagSyslog{},
			expectedError:           false,
			expectedFallthroughCall: true,
		},
		"error default, working fallthrough": {
			inputDefault:            &errorSyslog{},
			inputFallthrough:        &flagSyslog{},
			expectedError:           false,
			expectedFallthroughCall: true,
		},
		"working default, working fallthrough": {
			inputDefault:            &flagSyslog{},
			inputFallthrough:        &flagSyslog{},
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

		if fs, ok := test.inputDefault.(*flagSyslog); ok {
			actualCall := fs.flag

			assert.Equal(
				t,
				actualCall,
				test.expectedDefaultCall,
				fmt.Sprintf(
					"Fallthrough test call check failure"+
						" on default syslogger for: %s",
					explanation,
				),
			)
		}

		if fs, ok := test.inputFallthrough.(*flagSyslog); ok {
			actualCall := fs.flag

			assert.Equal(
				t,
				actualCall,
				test.expectedFallthroughCall,
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

package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultiSyslog(t *testing.T) {
	type testCase struct {
		inputSysloggers          []Syslogger
		inputTryAll              bool
		expectedError            bool
		expectedErrorSourceIndex int
		expectedCall             []bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputSysloggers: nil,
			expectedError:   false,
		},
		"one valid": {
			inputSysloggers: []Syslogger{
				&flagSyslog{},
			},
			expectedError: false,
			expectedCall: []bool{
				true,
			},
		},
		"one error": {
			inputSysloggers: []Syslogger{
				&errorSyslog{},
			},
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"one error, try all": {
			inputSysloggers: []Syslogger{
				&errorSyslog{},
			},
			inputTryAll:              true,
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"three valid": {
			inputSysloggers: []Syslogger{
				&flagSyslog{},
				&flagSyslog{},
				&flagSyslog{},
			},
			expectedError: false,
			expectedCall: []bool{
				true,
				true,
				true,
			},
		},
		"middle error": {
			inputSysloggers: []Syslogger{
				&flagSyslog{},
				&errorSyslog{},
				&flagSyslog{},
			},
			expectedError:            true,
			expectedErrorSourceIndex: 1,
			expectedCall: []bool{
				true,
				false,
				false,
			},
		},
		"middle error, try all": {
			inputSysloggers: []Syslogger{
				&flagSyslog{},
				&errorSyslog{},
				&flagSyslog{},
			},
			inputTryAll:              true,
			expectedError:            true,
			expectedErrorSourceIndex: 1,
			expectedCall: []bool{
				true,
				false,
				true,
			},
		},
		"two errors": {
			inputSysloggers: []Syslogger{
				&errorSyslog{},
				&errorSyslog{},
			},
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"two errors, try all": {
			inputSysloggers: []Syslogger{
				&errorSyslog{},
				&errorSyslog{},
			},
			inputTryAll:              true,
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"last two errors": {
			inputSysloggers: []Syslogger{
				&flagSyslog{},
				&errorSyslog{},
				&errorSyslog{},
			},
			expectedError:            true,
			expectedErrorSourceIndex: 1,
			expectedCall: []bool{
				true,
				false,
				false,
			},
		},
		"last errors, try all": {
			inputSysloggers: []Syslogger{
				&flagSyslog{},
				&errorSyslog{},
				&errorSyslog{},
			},
			inputTryAll:              true,
			expectedError:            true,
			expectedErrorSourceIndex: 1,
			expectedCall: []bool{
				true,
				false,
				false,
			},
		},
	}

	for explanation, test := range tests {
		m := &Multi{
			Sysloggers: test.inputSysloggers,
			TryAll:     test.inputTryAll,
		}

		actualError := m.Syslog(pri.Priority(0x0), nil)

		if test.expectedError {
			assert.Error(
				t,
				actualError,
				fmt.Sprintf(
					"Multi test expected error for: %s",
					explanation,
				),
			)

			if test.expectedErrorSourceIndex != 0 {
				i := test.expectedErrorSourceIndex

				assert.IsType(
					t,
					&errorSyslogError{},
					actualError,
					fmt.Sprintf(
						"Multi test expected error"+
							" from a log/"+
							"syslogger.errorSyslog"+
							" for: %s",
						explanation,
					),
				)

				ese, ok := actualError.(*errorSyslogError)
				if ok {
					assert.Equal(
						t,
						test.inputSysloggers[i],
						ese.s,
						fmt.Sprintf(
							"Multi test expected"+
								" error to"+
								" come from"+
								" syslogger %d"+
								" for: %s",
							i,
							explanation,
						),
					)
				}
			}
		} else {
			assert.NoError(
				t,
				actualError,
				fmt.Sprintf(
					"Multi test unexpected error for: %s",
					explanation,
				),
			)
		}

		for idx, s := range test.inputSysloggers {
			if fs, ok := s.(*flagSyslog); ok {
				actualCall := fs.flag

				assert.Equal(
					t,
					actualCall,
					test.expectedCall[idx],
					fmt.Sprintf(
						"Multi test call check failure"+
							" on default syslogger"+
							" for: %s",
						explanation,
					),
				)
			}
		}
	}
}

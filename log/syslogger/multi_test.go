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
				&flagSyslogger{},
			},
			expectedError: false,
			expectedCall: []bool{
				true,
			},
		},
		"one error": {
			inputSysloggers: []Syslogger{
				&errorSyslogger{},
			},
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"one error, try all": {
			inputSysloggers: []Syslogger{
				&errorSyslogger{},
			},
			inputTryAll:              true,
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"three valid": {
			inputSysloggers: []Syslogger{
				&flagSyslogger{},
				&flagSyslogger{},
				&flagSyslogger{},
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
				&flagSyslogger{},
				&errorSyslogger{},
				&flagSyslogger{},
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
				&flagSyslogger{},
				&errorSyslogger{},
				&flagSyslogger{},
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
				&errorSyslogger{},
				&errorSyslogger{},
			},
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"two errors, try all": {
			inputSysloggers: []Syslogger{
				&errorSyslogger{},
				&errorSyslogger{},
			},
			inputTryAll:              true,
			expectedError:            true,
			expectedErrorSourceIndex: 0,
		},
		"last two errors": {
			inputSysloggers: []Syslogger{
				&flagSyslogger{},
				&errorSyslogger{},
				&errorSyslogger{},
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
				&flagSyslogger{},
				&errorSyslogger{},
				&errorSyslogger{},
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
				expectedType := &errorSysloggerError{}

				assert.IsType(
					t,
					expectedType,
					actualError,
					fmt.Sprintf(
						"Multi test expected error"+
							" from a log/"+
							"syslogger.errorSyslog"+
							" for: %s",
						explanation,
					),
				)

				ese, ok := actualError.(*errorSysloggerError)
				if ok {
					expectedSrc := test.inputSysloggers[i]
					actualSrc := ese.S

					assert.Equal(
						t,
						expectedSrc,
						actualSrc,
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
			if fs, ok := s.(*flagSyslogger); ok {
				actualCall := fs.Flag

				assert.Equal(
					t,
					test.expectedCall[idx],
					actualCall,
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

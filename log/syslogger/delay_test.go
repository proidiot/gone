package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelayNewDelay(t *testing.T) {
	type testCase struct {
		inputCb                        func() (Syslogger, error)
		expectedNewDelayError          bool
		expectedSyslogError            bool
		expectedSysloggerCreationError bool
		expectedSyslogger              Syslogger
	}

	el := errorSyslog{}
	fl := flagSyslog{}

	tests := map[string]testCase{
		"nil values": {
			inputCb:               nil,
			expectedNewDelayError: true,
		},
		"callback error": {
			inputCb: func() (Syslogger, error) {
				return nil, errors.New(
					"callback error",
				)
			},
			expectedNewDelayError:          false,
			expectedSyslogError:            true,
			expectedSysloggerCreationError: true,
		},
		"error syslogger": {
			inputCb: func() (Syslogger, error) {
				return &el, nil
			},
			expectedNewDelayError:          false,
			expectedSyslogError:            true,
			expectedSysloggerCreationError: false,
			expectedSyslogger:              &el,
		},
		"flag syslogger": {
			inputCb: func() (Syslogger, error) {
				return &fl, nil
			},
			expectedNewDelayError:          false,
			expectedSyslogError:            false,
			expectedSysloggerCreationError: false,
			expectedSyslogger:              &fl,
		},
	}

	for explanation, test := range tests {
		n, actualNewDelayError := NewDelay(test.inputCb)

		if test.expectedNewDelayError {
			assert.Error(
				t,
				actualNewDelayError,
				fmt.Sprintf(
					"Delay test expected error for: %s",
					explanation,
				),
			)
		} else {
			assert.NoError(
				t,
				actualNewDelayError,
				fmt.Sprintf(
					"Delay test unexpected error for: %s",
					explanation,
				),
			)
		}

		if n != nil {
			actualError := n.Syslog(pri.Priority(0x0), nil)

			if test.expectedSyslogError {
				assert.Error(
					t,
					actualError,
					fmt.Sprintf(
						"Delay test expected error"+
							" for: %s",
						explanation,
					),
				)
			} else {
				assert.NoError(
					t,
					actualError,
					fmt.Sprintf(
						"Delay test unexpected error"+
							" for: %s",
						explanation,
					),
				)
			}

			actualSyslogHandler := n.h

			if test.expectedSysloggerCreationError {
				assert.Nil(
					t,
					actualSyslogHandler,
					fmt.Sprintf(
						"Delay test expected the"+
							" syslogger handler to"+
							" be nil due to errors"+
							" during creation for:"+
							" %s",
						explanation,
					),
				)
			} else {
				assert.NotNil(
					t,
					actualSyslogHandler,
					fmt.Sprintf(
						"Delay test expected the"+
							" syslogger handler to"+
							" not be nil for: %s",
						explanation,
					),
				)
			}

			if actualSyslogHandler != nil {
				actualSyslogger := actualSyslogHandler.s

				assert.Equal(
					t,
					test.expectedSyslogger,
					actualSyslogger,
					fmt.Sprintf(
						"Delay test did not generate"+
							" the expected"+
							" syslogger for: %s",
						explanation,
					),
				)
			}
		}
	}
}

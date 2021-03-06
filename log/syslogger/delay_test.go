package syslogger

import (
	"testing"

	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
)

func TestDelayNewDelay(t *testing.T) {
	type testCase struct {
		inputCb                        func() (Syslogger, error)
		expectedNewDelayError          bool
		expectedSyslogError            bool
		expectedSysloggerCreationError bool
		expectedSyslogger              Syslogger
	}

	el := errorSyslogger{}
	fl := flagSyslogger{}

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
			assert.Errorf(
				t,
				actualNewDelayError,
				"Delay test expected error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualNewDelayError,
				"Delay test unexpected error for: %s",
				explanation,
			)
		}

		if n != nil {
			actualError := n.Syslog(pri.Priority(0x0), nil)

			if test.expectedSyslogError {
				assert.Errorf(
					t,
					actualError,
					"Delay test expected error for: %s",
					explanation,
				)
			} else {
				assert.NoError(
					t,
					actualError,
					"Delay test unexpected error for: %s",
					explanation,
				)
			}

			actualSyslogHandler := n.h

			if test.expectedSysloggerCreationError {
				assert.Nil(
					t,
					actualSyslogHandler,
					"Delay test expected the syslogger"+
						" handler to be nil due to"+
						" errors during creation for:"+
						" %s",
					explanation,
				)
			} else {
				assert.NotNil(
					t,
					actualSyslogHandler,
					"Delay test expected the syslogger"+
						" handler to not be nil for:"+
						" %s",
					explanation,
				)
			}

			if actualSyslogHandler != nil {
				actualSyslogger := actualSyslogHandler.s

				assert.Equal(
					t,
					test.expectedSyslogger,
					actualSyslogger,
					"Delay test did not generate the"+
						" expected syslogger for: %s",
					explanation,
				)
			}
		}
	}
}

func TestDelayThreadSafety(t *testing.T) {
	s, sc := newSyncCountSyslogger()
	defer close(sc)

	d, e := NewDelay(
		func() (Syslogger, error) {
			return s, nil
		},
	)
	assert.NoError(
		t,
		e,
		"Delay thread safety test does not expect an error during"+
			" syslogger.Delay creation",
	)

	ec := make(chan error)
	defer close(ec)

	testFunc := func(msg string) {
		ec <- d.Syslog(pri.Priority(0x0), msg)
	}

	go testFunc("1")
	go testFunc("2")

	expectedInitialCount := uint(0)
	actualInitialCount := s.Count
	assert.Equal(
		t,
		expectedInitialCount,
		actualInitialCount,
		"Delay thread safety test expects a different initial count",
	)

	sc <- nil

	e = <-ec
	assert.NoError(
		t,
		e,
		"Delay thread safety test does not expect an error during"+
			" the first syslog call",
	)

	expectedMiddleCount := uint(1)
	actualMiddleCount := s.Count
	assert.Equal(
		t,
		expectedMiddleCount,
		actualMiddleCount,
		"Delay thread safety test expects a different middle count",
	)

	sc <- nil

	e = <-ec
	assert.NoError(
		t,
		e,
		"Delay thread safety test does not expect an error during"+
			" the second syslog call",
	)

	expectedLastCount := uint(2)
	actualLastCount := s.Count
	assert.Equal(
		t,
		expectedLastCount,
		actualLastCount,
		"Delay thread safety test expects a different last count",
	)
}

package syslogger

import (
	"testing"

	"github.com/proidiot/gone/log/mask"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
)

func TestSeverityMaskSyslog(t *testing.T) {
	type testCase struct {
		inputMask      mask.Mask
		inputPri       pri.Priority
		inputBadSyslog bool
		expectedError  bool
		expectedCall   bool
	}

	tests := map[string]testCase{
		"zero values": {
			inputMask:    mask.Mask(0x0),
			inputPri:     pri.Priority(0x0),
			expectedCall: false,
		},
		"full values": {
			inputMask:     mask.Mask(0xFF),
			inputPri:      pri.Priority(0xFF),
			expectedError: false,
			expectedCall:  true,
		},
		"notice explicitly unmasked": {
			inputMask:     mask.Notice,
			inputPri:      pri.Notice,
			expectedError: false,
			expectedCall:  true,
		},
		"notice not explicitly unmasked": {
			inputMask:    mask.Info,
			inputPri:     pri.Notice,
			expectedCall: false,
		},
		"notice upto unmasked": {
			inputMask:     mask.UpTo(pri.Info),
			inputPri:      pri.Notice,
			expectedError: false,
			expectedCall:  true,
		},
	}

	for explanation, test := range tests {
		es := errorSyslogger{}
		fs := flagSyslogger{}

		var s2 Syslogger

		if test.inputBadSyslog {
			s2 = &es
		} else {
			s2 = &fs
		}

		s := &SeverityMask{
			Syslogger: s2,
			Mask:      test.inputMask,
		}

		actualError := s.Syslog(test.inputPri, nil)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"SeverityMask test expected error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"SeverityMask test unexpected error for: %s",
				explanation,
			)
		}

		if !test.inputBadSyslog {
			actualCall := fs.Flag

			assert.Equal(
				t,
				test.expectedCall,
				actualCall,
				"SeverityMask test call check failure for: %s",
				explanation,
			)
		}
	}
}

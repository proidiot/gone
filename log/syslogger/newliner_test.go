package syslogger

import (
	"testing"

	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
)

func TestNewlinerSyslog(tt *testing.T) {
	tests := map[string]struct {
		inputPri       pri.Priority
		inputMsg       interface{}
		expectedError  bool
		expectedOutput string
	}{
		"zero values": {
			inputPri:      pri.Priority(0x00),
			inputMsg:      []byte{},
			expectedError: true,
		},
		"full values": {
			inputPri:      pri.Priority(0xFF),
			inputMsg:      []byte{0xFF, 0xFF, 0xFF, 0xFF},
			expectedError: true,
		},
		"no newline": {
			inputMsg:       "testing",
			expectedOutput: "testing\n",
		},
		"one newline": {
			inputMsg:       "testing\n",
			expectedOutput: "testing\n",
		},
		"two newlines": {
			inputMsg:       "testing\n\n",
			expectedOutput: "testing\n\n",
		},
		"dos newlines": {
			inputMsg:       "testing\r\n",
			expectedOutput: "testing\r\n",
		},
		"reverse dos newlines": {
			inputMsg:       "testing\n\r",
			expectedOutput: "testing\n\r\n",
		},
		"almost dos newlines": {
			inputMsg:       "testing\r",
			expectedOutput: "testing\r\n",
		},
		"middle newline": {
			inputMsg:       "testing\n123",
			expectedOutput: "testing\n123\n",
		},
		"only newline": {
			inputMsg:       "\n",
			expectedOutput: "\n",
		},
		"empty string": {
			inputMsg:       "",
			expectedOutput: "\n",
		},
		"stringer": {
			inputMsg:       &stringer{"test"},
			expectedOutput: "test\n",
		},
	}

	for explanation, test := range tests {
		tt.Run(explanation, func(t *testing.T) {
			//t.Parallel()

			rec := new(recordStringSyslogger)

			newliner := &Newliner{
				Syslogger: rec,
			}

			actualError := newliner.Syslog(
				test.inputPri,
				test.inputMsg,
			)

			if test.expectedError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			actualOutput := rec.M
			assert.Equal(
				t,
				test.expectedOutput,
				actualOutput,
			)
		})
	}
}

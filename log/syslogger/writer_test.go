package syslogger

import (
	"bytes"
	"testing"

	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
)

func TestWriterSyslog(t *testing.T) {
	type testCase struct {
		inputPri       pri.Priority
		inputMsg       interface{}
		inputBadWriter bool
		expectedError  bool
		expectedBytes  []byte
	}

	tests := map[string]testCase{
		"zero values": {
			inputPri:      pri.Priority(0x00),
			inputMsg:      []byte{},
			expectedError: false,
			expectedBytes: nil,
		},
		"full values": {
			inputPri:      pri.Priority(0xFF),
			inputMsg:      []byte{0xFF, 0xFF, 0xFF, 0xFF},
			expectedError: true,
		},
		"string message": {
			inputMsg:      "testing",
			expectedBytes: []byte("testing"),
		},
		"bytes message": {
			inputMsg:      []byte{0xC0, 0xA8, 0x00, 0x01},
			expectedBytes: []byte{0xC0, 0xA8, 0x00, 0x01},
		},
		"struct message": {
			inputMsg:      testCase{},
			expectedError: true,
		},
		"string message bad writer": {
			inputMsg:       "testing",
			inputBadWriter: true,
			expectedError:  true,
		},
		"bytes message bad writer": {
			inputMsg:       []byte{0xC0, 0xA8, 0x00, 0x01},
			inputBadWriter: true,
			expectedError:  true,
		},
	}

	for explanation, test := range tests {
		actualBytes := new(bytes.Buffer)

		var writer *Writer
		if test.inputBadWriter {
			writer = &Writer{
				Writer: errorWriter{},
			}
		} else {
			writer = &Writer{
				Writer: actualBytes,
			}
		}

		actualError := writer.Syslog(
			test.inputPri,
			test.inputMsg,
		)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"Writer test expected error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"Writer test unexpected error for: %s",
				explanation,
			)
		}

		assert.Equal(
			t,
			test.expectedBytes,
			actualBytes.Bytes(),
			"Writer test bytes check failure for: %s",
			explanation,
		)
	}
}

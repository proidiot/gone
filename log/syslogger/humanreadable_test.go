package syslogger

import (
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestHumanReadableSyslog(t *testing.T) {
	dateregex := `.*`
	hostregex := `[0-9A-Za-z][-0-9A-Za-z]*[0-9A-Za-z]`

	type testCase struct {
		inputFacility        pri.Priority
		inputIdent           string
		inputPid             bool
		inputPiority         pri.Priority
		inputMsg             interface{}
		causeOsHostnameError bool
		expectedError        bool
		expectedPriority     pri.Priority
		expectedMsg          *regexp.Regexp
	}

	tests := map[string]testCase{
		"nil values": {
			inputFacility: pri.Priority(0x0),
			inputIdent:    string(0),
			inputPid:      false,
			inputPiority:  pri.Priority(0x0),
			inputMsg:      nil,
			expectedError: true,
		},
		"full values": {
			inputFacility:    pri.Priority(0xFF),
			inputIdent:       "full values ident",
			inputPid:         true,
			inputPiority:     pri.Priority(0xFF),
			inputMsg:         "full values message",
			expectedError:    false,
			expectedPriority: pri.Priority(0x0),
			expectedMsg: regexp.MustCompile(
				`^Priority\(0xf8\) LOG_DEBUG ` + dateregex +
					` ` + hostregex + ` full values ident` +
					`\[\d+\] full values message$`,
			),
		},
		"normal call": {
			inputPiority:  pri.Notice,
			inputMsg:      "normal call message",
			expectedError: false,
			expectedMsg: regexp.MustCompile(
				`LOG_USER LOG_NOTICE ` + dateregex + ` ` +
					hostregex +
					` [^ ]+ normal call message$`,
			),
		},
		"error call": {
			inputPiority:  pri.Err,
			inputMsg:      errors.New("error call message"),
			expectedError: false,
			expectedMsg: regexp.MustCompile(
				`LOG_USER LOG_ERR ` + dateregex + ` ` +
					hostregex +
					` [^ ]+ error call message$`,
			),
		},
		"stringer call": {
			inputFacility: pri.Syslog,
			inputPiority:  pri.Debug,
			inputMsg:      pri.Debug,
			expectedError: false,
			expectedMsg: regexp.MustCompile(
				`LOG_SYSLOG LOG_DEBUG ` + dateregex + ` ` +
					hostregex +
					` [^ ]+ LOG_DEBUG`,
			),
		},
		"broken hostname": {
			inputPiority:         pri.Info,
			inputMsg:             "broken hostname message",
			causeOsHostnameError: true,
			expectedError:        false,
			expectedMsg: regexp.MustCompile(
				`LOG_USER LOG_INFO ` + dateregex +
					` localhost [^ ]+` +
					` broken hostname message`,
			),
		},
	}

	for explanation, test := range tests {
		origOsHostname := osHostname
		if test.causeOsHostnameError {
			osHostname = func() (string, error) {
				return "", errors.New(
					"Artificial error for os.Hostname",
				)
			}
			defer func() {
				osHostname = origOsHostname
			}()
		}

		rs := recordStringSyslogger{}

		h := &HumanReadable{
			Syslogger: &rs,
			Facility:  test.inputFacility,
			Ident:     test.inputIdent,
			Pid:       test.inputPid,
		}

		actualError := h.Syslog(test.inputPiority, test.inputMsg)

		if test.expectedError {
			assert.Error(
				t,
				actualError,
				"HumanReadable test expects an error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"HumanReadable test expects no error for: %s",
				explanation,
			)
		}

		actualPriority := rs.P
		assert.Equal(
			t,
			test.expectedPriority,
			actualPriority,
			"HumanReadable test recorded the wrong pri.Priority"+
				" for: %s",
			explanation,
		)

		if test.expectedMsg != nil {
			actualMsg := rs.M
			assert.Regexp(
				t,
				test.expectedMsg,
				actualMsg,
				"HumanReadable test recorded a non-matching"+
					" string for: %s",
				explanation,
			)
		}

		osHostname = origOsHostname
	}
}

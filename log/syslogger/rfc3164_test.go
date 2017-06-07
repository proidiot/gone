package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRfc3164Syslog(t *testing.T) {
	dateregex := `((Jan)|(Feb)|(Mar)|(Apr)|(May)|(Jun)|(Jul)|(Aug)|(Sep)` +
		`|(Oct)|(Nov)|(Dec)) ((30)|(31)|([ 12]\d))` +
		` ((2[0-3])|([ 01]\d)):([0-5]\d):([0-5]\d)`
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
				`^<255>` + dateregex + ` ` + hostregex +
					` full values ident\[\d+\]:` +
					` full values message$`,
			),
		},
		"normal call": {
			inputPiority:  pri.Notice,
			inputMsg:      "normal call message",
			expectedError: false,
			expectedMsg: regexp.MustCompile(
				`<13>` + dateregex + ` ` + hostregex +
					` [^ ]+: normal call message$`,
			),
		},
		"error call": {
			inputPiority:  pri.Err,
			inputMsg:      errors.New("error call message"),
			expectedError: true,
		},
		"stringer call": {
			inputFacility: pri.Syslog,
			inputPiority:  pri.Debug,
			inputMsg:      pri.Debug,
			expectedError: true,
		},
		"broken hostname": {
			inputPiority:         pri.Info,
			inputMsg:             "broken hostname message",
			causeOsHostnameError: true,
			expectedError:        false,
			expectedMsg: regexp.MustCompile(
				`<14>` + dateregex +
					` localhost [^ ]+` +
					` broken hostname message`,
			),
		},
		"ridiculously long message": {
			inputPiority: pri.Warning,
			inputMsg: "This message is way too long to be used as" +
				" the message to be specified within an" +
				" actual syslog call because the Internet" +
				" Engineering Task Force's Request for" +
				" Comments number three thousand one hundred" +
				" and sixty-four specifies that conforming" +
				" messages must have no more than one" +
				" thousand and twenty-four characters (with" +
				" such character being limited to those" +
				" specified within the American Standard Code" +
				" for Information Interchange), and yet this" +
				" very message contains more than one" +
				" thousand and twenty-four such characters." +
				" Specifically, the Internet Engineering Task" +
				" Force's Request for Comments number three" +
				" thousand one hundred and sixty-four" +
				" references this particular restriction in" +
				" the following locations: the leading" +
				" paragraph of section four, subsection one;" +
				" the final paragraph of section four," +
				" subsection three, sub-subsection two; the" +
				" final paragraph of section four, subsection" +
				" three, sub-subsection three; and, most" +
				" prominently, in section six subsection one," +
				" where there is presented a more complete" +
				" explanation of why this particular" +
				" restriction has been imposed on all" +
				" messages wishing to conform to this" +
				" standard. In fact, this message has one" +
				" thousand two hundred and one characters.",
			expectedError: true,
		},
	}

	for explanation, test := range tests {
		var origOsHostname func() (string, error) = osHostname
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

		rs := recordStringSyslog{}

		h := &Rfc3164{
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
				fmt.Sprintf(
					"Rfc3164 test expects an error"+
						" for: %s",
					explanation,
				),
			)
		} else {
			assert.NoError(
				t,
				actualError,
				fmt.Sprintf(
					"Rfc3164 test expects no error"+
						" for: %s",
					explanation,
				),
			)
		}

		actualPriority := rs.p
		assert.Equal(
			t,
			test.expectedPriority,
			actualPriority,
			fmt.Sprintf(
				"Rfc3164 test recorded the wrong"+
					" pri.Priority for: %s",
				explanation,
			),
		)

		if test.expectedMsg != nil {
			actualMsg := rs.m
			assert.Regexp(
				t,
				test.expectedMsg,
				actualMsg,
				fmt.Sprintf(
					"Rfc3164 test recorded a"+
						" non-matching string for: %s",
					explanation,
				),
			)
		}

		osHostname = origOsHostname
	}
}

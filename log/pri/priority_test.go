package pri

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriFacility(t *testing.T) {
	type testCase struct {
		input    Priority
		expected Priority
	}

	tests := map[string]testCase{
		"zero priority": {
			input:    Priority(0),
			expected: Priority(0),
		},
		"user facility": {
			input:    User,
			expected: User,
		},
		"err severity": {
			input:    Err,
			expected: Priority(0),
		},
		"news warning combo": {
			input:    News | Warning,
			expected: News,
		},
		"nonsensical multi facility": {
			input:    Local0 | Uucp,
			expected: Priority(0xC0),
		},
		"full byte": {
			input:    0xFF,
			expected: Priority(0xF8),
		},
	}

	for explanation, test := range tests {
		actual := test.input.Facility()

		assert.Equal(
			t,
			test.expected,
			actual,
			"Facility test failed for test case: %s",
			explanation,
		)
	}
}

func TestPriValidFacility(t *testing.T) {
	type testCase struct {
		input         Priority
		errorExpected bool
	}

	tests := map[string]testCase{
		"zero priority": {
			input:         Priority(0),
			errorExpected: false,
		},
		"user facility": {
			input:         User,
			errorExpected: false,
		},
		"err severity": {
			input:         Err,
			errorExpected: true,
		},
		"news warning combo": {
			input:         News | Warning,
			errorExpected: true,
		},
		"nonsensical multi facility": {
			input:         Local0 | Uucp,
			errorExpected: true,
		},
		"full byte": {
			input:         0xFF,
			errorExpected: true,
		},
	}

	for explanation, test := range tests {
		actualError := test.input.ValidFacility()

		if test.errorExpected {
			assert.Errorf(
				t,
				actualError,
				"ValidFacility test had false negative for"+
					" test case: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"ValidFacility test had false positive for"+
					" test case: %s",
				explanation,
			)
		}
	}
}

func TestPriSeverity(t *testing.T) {
	type testCase struct {
		input    Priority
		expected Priority
	}

	tests := map[string]testCase{
		"zero priority": {
			input:    Priority(0),
			expected: Priority(0),
		},
		"user facility": {
			input:    User,
			expected: Priority(0),
		},
		"err severity": {
			input:    Err,
			expected: Err,
		},
		"news warning combo": {
			input:    News | Warning,
			expected: Warning,
		},
		"nonsensical multi facility": {
			input:    Local0 | Uucp,
			expected: Priority(0),
		},
		"full byte": {
			input:    0xFF,
			expected: Debug,
		},
	}

	for explanation, test := range tests {
		actual := test.input.Severity()

		assert.Equal(
			t,
			test.expected,
			actual,
			"Severity test failed for test case: %s",
			explanation,
		)
	}
}

func TestPriString(t *testing.T) {
	type testCase struct {
		input    Priority
		expected string
	}

	tests := map[string]testCase{
		"zero priority": {
			input:    Priority(0),
			expected: "LOG_EMERG",
		},
		"user facility": {
			input:    User,
			expected: "LOG_USER",
		},
		"err severity": {
			input:    Err,
			expected: "LOG_ERR",
		},
		"news warning combo": {
			input:    News | Warning,
			expected: "LOG_NEWS|LOG_WARNING",
		},
		"nonsensical multi facility": {
			input:    Local0 | Uucp,
			expected: "Priority(0xc0)",
		},
		"full byte": {
			input:    0xFF,
			expected: "Priority(0xf8)|LOG_DEBUG",
		},
	}

	for explanation, test := range tests {
		actual := test.input.String()

		assert.Equal(
			t,
			test.expected,
			actual,
			"String test failed for test case: %s",
			explanation,
		)
	}
}

func TestPriGetFromEnv(t *testing.T) {
	clearEnvs := []string{
		"LOG_FACILITY",
		"LOG_PRIORITY",
	}

	type testCase struct {
		input    map[string]string
		expected Priority
	}

	tests := map[string]testCase{
		"no vals set": {
			input:    map[string]string{},
			expected: User,
		},
		"facility set": {
			input: map[string]string{
				"LOG_FACILITY": "LOG_CRON",
			},
			expected: Cron,
		},
		"priority set": {
			input: map[string]string{
				"LOG_PRIORITY": "LOG_FTP",
			},
			expected: Ftp,
		},
		"facility overrides priority": {
			input: map[string]string{
				"LOG_PRIORITY": "LOG_NTP",
				"LOG_FACILITY": "LOG_LPR",
			},
			expected: Lpr,
		},
		"severity as facility": {
			input: map[string]string{
				"LOG_FACILITY": "LOG_INFO",
			},
			expected: User,
		},
		"combo as facility": {
			input: map[string]string{
				"LOG_FACILITY": "LOG_LOCAL0|LOG_AUDIT",
			},
			expected: User,
		},
		"numerical facility": {
			input: map[string]string{
				"LOG_FACILITY": "40",
			},
			expected: User, // TODO should this change?
		},
	}

	for explanation, test := range tests {
		for _, env := range clearEnvs {
			e := os.Unsetenv(env)
			assert.NoError(t, e, "Error during Unsetenv")
		}

		for env, val := range test.input {
			e := os.Setenv(env, val)
			assert.NoError(t, e, "Error during Setenv")
		}

		actual := GetFromEnv()

		assert.Equal(
			t,
			test.expected,
			actual,
			"GetFromEnv test failed for: %s",
			explanation,
		)
	}
}

package mask

import (
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMaskUpTo(t *testing.T) {
	type testCase struct {
		input    pri.Priority
		expected Mask
	}

	tests := map[string]testCase{
		"zero value input gives non-zero value outpu": {
			input:    pri.Priority(0x0),
			expected: Mask(0x01),
		},
		"max input gives max output": {
			input:    pri.Priority(0xFF),
			expected: Mask(0xFF),
		},
		"specific input gives specific output": {
			input:    pri.Err,
			expected: Err | Crit | Alert | Emerg,
		},
		"facility does not effect output": {
			input:    pri.Ntp | pri.Crit,
			expected: Crit | Alert | Emerg,
		},
	}

	for explanation, test := range tests {
		actual := UpTo(test.input)

		assert.Equal(
			t,
			test.expected,
			actual,
			"UpTo test failed for: %s",
			explanation,
		)
	}
}

func TestMaskMasked(t *testing.T) {
	type testCase struct {
		inputMask Mask
		inputPri  pri.Priority
		expected  bool
	}

	tests := map[string]testCase{
		"zero mask hides everything": {
			inputMask: Mask(0x0),
			inputPri:  pri.Priority(0x0),
			expected:  true,
		},
		"full mask hides nothing": {
			inputMask: Mask(0xFF),
			inputPri:  pri.Priority(0xFF),
			expected:  false,
		},
		"specific mask allows specific level": {
			inputMask: Warning,
			inputPri:  pri.Warning,
			expected:  false,
		},
		"mask doesn't imply lower masks": {
			inputMask: Crit,
			inputPri:  pri.Alert,
			expected:  true,
		},
		"bitwise or of mask doesn't imply bitwise or of priority": {
			inputMask: Alert | Crit,
			inputPri:  pri.Err,
			expected:  true,
		},
		"compound mask allows constituents": {
			inputMask: Err | Warning | Notice,
			inputPri:  pri.Warning,
			expected:  false,
		},
	}

	for explanation, test := range tests {
		actual := test.inputMask.Masked(test.inputPri)

		assert.Equal(
			t,
			test.expected,
			actual,
			"Masked test failed for: %s",
			explanation,
		)
	}
}

func TestMaskString(t *testing.T) {
	type testCase struct {
		input    Mask
		expected string
	}

	tests := map[string]testCase{
		"zero value": {
			input:    Mask(0x00),
			expected: "LOG_MASK(0x0)",
		},
		"full value": {
			input:    Mask(0xFF),
			expected: "LOG_UPTO(LOG_DEBUG)",
		},
		"specific value": {
			input:    Err,
			expected: "LOG_MASK(LOG_ERR)",
		},
		"upto value": {
			input:    Crit | Alert | Emerg,
			expected: "LOG_UPTO(LOG_CRIT)",
		},
		"multi value": {
			input:    Err | Alert | Emerg,
			expected: "LOG_MASK(LOG_EMERG|LOG_ALERT|LOG_ERR)",
		},
	}

	for explanation, test := range tests {
		actual := test.input.String()

		assert.Equal(
			t,
			test.expected,
			actual,
			"String test failed for: %s",
			explanation,
		)
	}
}

func TestMaskGetFromEnv(t *testing.T) {
	clearEnvs := []string{
		"LOG_UPTO",
		"LOG_MASK",
	}

	type testCase struct {
		input    map[string]string
		expected Mask
	}

	tests := map[string]testCase{
		"no vals set": {
			input:    map[string]string{},
			expected: Mask(0xFF),
		},
		"upto set": {
			input: map[string]string{
				"LOG_UPTO": "LOG_CRIT",
			},
			expected: Crit | Alert | Emerg,
		},
		"mask set": {
			input: map[string]string{
				"LOG_MASK": "LOG_INFO",
			},
			expected: Info,
		},
		"both set": {
			input: map[string]string{
				"LOG_UPTO": "LOG_ALERT",
				"LOG_MASK": "LOG_WARNING",
			},
			expected: Alert | Emerg,
		},
		"multi mask set": {
			input: map[string]string{
				"LOG_MASK": "LOG_ERR|LOG_NOTICE",
			},
			expected: Err | Notice,
		},
		"upto debug": {
			input: map[string]string{
				"LOG_UPTO": "LOG_DEBUG",
			},
			expected: Mask(0xFF),
		},
		"bad simple mask": {
			input: map[string]string{
				"LOG_MASK": "LOG_USER",
			},
			expected: Mask(0xFF),
		},
		"bad compound syntax": {
			input: map[string]string{
				"LOG_MASK": "LOG_ERR | LOG_NOTICE",
			},
			expected: Mask(0xFF),
		},
		"bad compound value": {
			input: map[string]string{
				"LOG_MASK": "LOG_Notice|LOG_KERN",
			},
			expected: Mask(0xFF),
		},
		"bad upto": {
			input: map[string]string{
				"LOG_UPTO": "LOG_ERROR",
			},
			expected: Mask(0xFF),
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

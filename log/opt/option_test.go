package opt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestOptGetFromEnv(t *testing.T) {
	envOpts := []string{
		"LOG_PID",
		"LOG_CONS",
		"LOG_ODELAY",
		"LOG_NDELAY",
		"LOG_NOWAIT",
		"LOG_PERROR",
		"LOG_NOFALLBACK",
	}

	type testCase struct{
		env []string
		expected Option
		explanation string
	}

	tests := []testCase{
		testCase{
			env: []string{},
			expected: 0x00,
			explanation: "no options set",
		},
		testCase{
			env: []string{
				"LOG_PID",
			},
			expected: Pid,
			explanation: "pid only",
		},
		testCase{
			env: []string{
				"LOG_CONS",
			},
			expected: Cons,
			explanation: "cons only",
		},
		testCase{
			env: []string{
				"LOG_ODELAY",
			},
			expected: ODelay,
			explanation: "odelay only",
		},
		testCase{
			env: []string{
				"LOG_NDELAY",
			},
			expected: NDelay,
			explanation: "ndelay only",
		},
		testCase{
			env: []string{
				"LOG_NOWAIT",
			},
			expected: NoWait,
			explanation: "nowait only",
		},
		testCase{
			env: []string{
				"LOG_PERROR",
			},
			expected: Perror,
			explanation: "perror only",
		},
		testCase{
			env: []string{
				"LOG_NOFALLBACK",
			},
			expected: NoFallback,
			explanation: "nofallback only",
		},
		testCase{
			env: []string{
				"LOG_PID",
				"LOG_CONS",
			},
			expected: Pid|Cons,
			explanation: "pid and cons",
		},
		testCase{
			env: []string{
				"LOG_PID",
				"LOG_CONS",
				"LOG_ODELAY",
				"LOG_NDELAY",
				"LOG_NOWAIT",
				"LOG_PERROR",
				"LOG_NOFALLBACK",
			},
			expected: 0x7F,
			explanation: "all opts",
		},
	}

	for _, test := range tests {
		for _, envOpt := range envOpts {
			e := os.Unsetenv(envOpt)
			assert.NoError(t, e, "Error during Unsetenv")
		}

		for _, env := range test.env {
			e := os.Setenv(env, "1")
			assert.NoError(t, e, "Error during Setenv")
		}

		actual := GetFromEnv()

		assert.Equal(
			t,
			test.expected,
			actual,
			fmt.Sprintf(
				"GetFromEnv test failed for: %s",
				test.explanation,
			),
		)
	}
}

func TestOptString(t *testing.T) {
	envOpts := []string{
		"LOG_PID",
		"LOG_CONS",
		"LOG_ODELAY",
		"LOG_NDELAY",
		"LOG_NOWAIT",
		"LOG_PERROR",
		"LOG_NOFALLBACK",
	}

	type testCase struct{
		input Option
		expected string
		explanation string
	}

	tests := []testCase{
		testCase{
			input: 0x00,
			expected: "Option(0)",
			explanation: "no options set",
		},
		testCase{
			input: Pid,
			expected: "LOG_PID",
			explanation: "pid only",
		},
		testCase{
			input: Cons,
			expected: "LOG_CONS",
			explanation: "cons only",
		},
		testCase{
			input: ODelay,
			expected: "LOG_ODELAY",
			explanation: "odelay only",
		},
		testCase{
			input: NDelay,
			expected: "LOG_NDELAY",
			explanation: "ndelay only",
		},
		testCase{
			input: NoWait,
			expected: "LOG_NOWAIT",
			explanation: "nowait only",
		},
		testCase{
			input: Perror,
			expected: "LOG_PERROR",
			explanation: "perror only",
		},
		testCase{
			input: NoFallback,
			expected: "LOG_NOFALLBACK",
			explanation: "nofallback only",
		},
		testCase{
			input: Pid|Cons,
			expected: "LOG_PID|LOG_CONS",
			explanation: "pid and cons",
		},
		testCase{
			input: Perror|NDelay|NoWait,
			expected: "LOG_NDELAY|LOG_NOWAIT|LOG_PERROR",
			explanation: "ndelay, nowait, and perror",
		},
		testCase{
			input: 0x7F,
			expected: "LOG_PID|LOG_CONS|LOG_ODELAY|LOG_NDELAY" +
				"|LOG_NOWAIT|LOG_PERROR|LOG_NOFALLBACK",
			explanation: "all opts",
		},
		testCase{
			input: 0x80,
			expected: "Option(80)",
			explanation: "unknown opt",
		},
		testCase{
			input: 0x84,
			expected: "LOG_ODELAY|Option(80)",
			explanation: "odelay and unknown opt",
		},
	}

	for _, test := range tests {
		for _, envOpt := range envOpts {
			e := os.Unsetenv(envOpt)
			assert.NoError(t, e, "Error during Unsetenv")
		}

		actual := test.input.String()

		assert.Equal(
			t,
			test.expected,
			actual,
			fmt.Sprintf(
				"GetFromEnv test failed for: %s",
				test.explanation,
			),
		)
	}
}

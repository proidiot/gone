package syslogger

import (
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/syslog"
	"net"
	"testing"
)

func TestNewNativeSyslog(t *testing.T) {
	var origSyslogNew = syslogNew
	defer func() {
		syslogNew = origSyslogNew
	}()

	tryRealSyslog := true
	if _, e := syslogNew(0, ""); e != nil {
		t.Log(
			"It seems that log/syslog.New is failing by default." +
				" This could be because the system doesn't" +
				" have a running syslogd (such as in a" +
				" container), or it could be because the" +
				" system isn't Unix-like. In order to" +
				" complete the tests, a mock of" +
				" log/syslog.New will be used for the happy" +
				" path in addition to the mock already" +
				" being used for the sad path.",
		)
		tryRealSyslog = false
	}

	type testCase struct {
		inputPriority       pri.Priority
		inputIdent          string
		causeSyslogNewError bool
		expectedError       bool
		expectedSyslogger   bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputPriority:     pri.Priority(0x0),
			inputIdent:        "",
			expectedError:     false,
			expectedSyslogger: true,
		},
		"full values": {
			inputPriority: pri.Priority(0xFF),
			inputIdent:    "full values ident",
			expectedError: true,
		},
		"broken syslog.New": {
			causeSyslogNewError: true,
			expectedError:       true,
		},
	}

	for explanation, test := range tests {
		if test.causeSyslogNewError {
			syslogNew = func(
				syslog.Priority,
				string,
			) (*syslog.Writer, error) {
				return nil, errors.New(
					"Artificial error for syslog.New",
				)
			}
		} else if tryRealSyslog {
			syslogNew = origSyslogNew
		} else {
			syslogNew = func(
				syslog.Priority,
				string,
			) (*syslog.Writer, error) {
				return &syslog.Writer{}, nil
			}
		}

		actualSyslogger, actualError := NewNativeSyslog(
			test.inputPriority,
			test.inputIdent,
		)

		if test.expectedError {
			assert.Error(
				t,
				actualError,
				fmt.Sprintf(
					"NewNativeSyslog test expects an"+
						" error for: %s",
					explanation,
				),
			)
		} else {
			assert.NoError(
				t,
				actualError,
				fmt.Sprintf(
					"NewNativeSyslog test expects no"+
						" error for: %s",
					explanation,
				),
			)
		}

		if test.expectedSyslogger {
			assert.NotNil(
				t,
				actualSyslogger,
				fmt.Sprintf(
					"NewNativeSyslog test expects non-nil"+
						" syslogger for: %s",
					explanation,
				),
			)
		} else {
			assert.Nil(
				t,
				actualSyslogger,
				fmt.Sprintf(
					"NewNativeSyslog test expects nil"+
						" syslogger for: %s",
					explanation,
				),
			)
		}

		if actualSyslogger != nil {
			e := actualSyslogger.Close()
			assert.NoError(
				t,
				e,
				fmt.Sprintf(
					"NewNativeSyslog test expects no"+
						" error when closing the"+
						" tested NativeSyslog for: %s",
					explanation,
				),
			)
		}
	}
}

func TestDialNativeSyslog(t *testing.T) {
	var origSyslogNew = syslogNew
	defer func() {
		syslogNew = origSyslogNew
	}()

	if _, e := syslogNew(0, ""); e != nil {
		t.Log(
			"It seems that log/syslog.New is failing by default." +
				" This could be because the system doesn't" +
				" have a running syslogd (such as in a" +
				" container), or it could be because the" +
				" system isn't Unix-like. In order to" +
				" complete the tests, a mock of" +
				" log/syslog.New will be used for the happy" +
				" path in addition to the mock already" +
				" being used for the sad path.",
		)
		syslogNew = func(
			syslog.Priority,
			string,
		) (*syslog.Writer, error) {
			return &syslog.Writer{}, nil
		}
	}

	udpNetwork := "udp"
	udpListener, e := net.ListenUDP(udpNetwork, &net.UDPAddr{})
	require.NoError(
		t,
		e,
		"DialNativeSyslog test expects no error when creating the"+
			" default udp listener.",
	)
	defer udpListener.Close()
	udpRaddr := udpListener.LocalAddr().String()

	tcpNetwork := "tcp"
	tcpListener, e := net.Listen(tcpNetwork, ":0")
	require.NoError(
		t,
		e,
		"DialNativeSyslog test expects no error when creating the"+
			" default tcp listener.",
	)
	defer tcpListener.Close()
	tcpRaddr := tcpListener.Addr().String()

	unixNetwork := "unix"
	unixListener, e := net.Listen(unixNetwork, "")
	require.NoError(
		t,
		e,
		"DialNativeSyslog test expects no error when creating the"+
			" default unix listener.",
	)
	defer unixListener.Close()
	unixRaddr := unixListener.Addr().String()

	badNetwork := "tcp4"
	badRaddr := "0.0.0.0:514"

	type testCase struct {
		inputNetwork      string
		inputRaddr        string
		inputPriority     pri.Priority
		inputIdent        string
		expectedError     bool
		expectedSyslogger bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputNetwork:      "",
			inputRaddr:        "",
			inputPriority:     pri.Priority(0x0),
			inputIdent:        "",
			expectedError:     false,
			expectedSyslogger: true,
		},
		"full values": {
			inputNetwork:  udpNetwork,
			inputRaddr:    udpRaddr,
			inputPriority: pri.Priority(0xFF),
			inputIdent:    "full values ident",
			expectedError: true,
		},
		"udp connection": {
			inputNetwork:      udpNetwork,
			inputRaddr:        udpRaddr,
			expectedError:     false,
			expectedSyslogger: true,
		},
		"tcp connection": {
			inputNetwork:      tcpNetwork,
			inputRaddr:        tcpRaddr,
			expectedError:     false,
			expectedSyslogger: true,
		},
		"unix connection": {
			inputNetwork:      unixNetwork,
			inputRaddr:        unixRaddr,
			expectedError:     false,
			expectedSyslogger: true,
		},
		"garbage connection": {
			inputNetwork:  "not a real network",
			inputRaddr:    "not a real address",
			expectedError: true,
		},
		"bad connection": {
			inputNetwork:  badNetwork,
			inputRaddr:    badRaddr,
			expectedError: true,
		},
	}

	for explanation, test := range tests {
		actualSyslogger, actualError := DialNativeSyslog(
			test.inputNetwork,
			test.inputRaddr,
			test.inputPriority,
			test.inputIdent,
		)

		if test.expectedError {
			assert.Error(
				t,
				actualError,
				fmt.Sprintf(
					"DialNativeSyslog test expects an"+
						" error for: %s",
					explanation,
				),
			)
		} else {
			assert.NoError(
				t,
				actualError,
				fmt.Sprintf(
					"DialNativeSyslog test expects no"+
						" error for: %s",
					explanation,
				),
			)
		}

		if test.expectedSyslogger {
			assert.NotNil(
				t,
				actualSyslogger,
				fmt.Sprintf(
					"DialNativeSyslog test expects non-nil"+
						" syslogger for: %s",
					explanation,
				),
			)
		} else {
			assert.Nil(
				t,
				actualSyslogger,
				fmt.Sprintf(
					"DialNativeSyslog test expects nil"+
						" syslogger for: %s",
					explanation,
				),
			)
		}

		if actualSyslogger != nil {
			e := actualSyslogger.Close()
			assert.NoError(
				t,
				e,
				fmt.Sprintf(
					"NewNativeSyslog test expects no"+
						" error when closing the"+
						" tested NativeSyslog for: %s",
					explanation,
				),
			)
		}
	}
}

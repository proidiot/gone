package syslogger

import (
	"bufio"
	"fmt"
	"github.com/proidiot/gone/errors"
	"github.com/proidiot/gone/log/pri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/syslog"
	"net"
	"regexp"
	"testing"
)

func TestNewNativeSyslog(t *testing.T) {
	origSyslogNew := syslogNew
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
			assert.Errorf(
				t,
				actualError,
				"NewNativeSyslog test expects an error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"NewNativeSyslog test expects no error for: %s",
				explanation,
			)
		}

		if test.expectedSyslogger {
			assert.NotNil(
				t,
				actualSyslogger,
				"NewNativeSyslog test expects non-nil"+
					" syslogger for: %s",
				explanation,
			)
		} else {
			assert.Nil(
				t,
				actualSyslogger,
				"NewNativeSyslog test expects nil syslogger"+
					" for: %s",
				explanation,
			)
		}

		if actualSyslogger != nil {
			e := actualSyslogger.Close()
			assert.NoError(
				t,
				e,
				"NewNativeSyslog test expects no error when"+
					" closing the tested NativeSyslog for:"+
					" %s",
				explanation,
			)
		}
	}
}

func TestDialNativeSyslog(t *testing.T) {
	syslogNewWorks := true
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
		syslogNewWorks = false
	}

	udpNetwork := "udp"
	udpListener, e := net.ListenUDP(udpNetwork, &net.UDPAddr{})
	require.NoError(
		t,
		e,
		"DialNativeSyslog test expects no error when creating the"+
			" default udp listener.",
	)
	defer func() {
		_ = udpListener.Close()
	}()
	udpRaddr := udpListener.LocalAddr().String()

	tcpNetwork := "tcp"
	tcpListener, e := net.Listen(tcpNetwork, ":0")
	require.NoError(
		t,
		e,
		"DialNativeSyslog test expects no error when creating the"+
			" default tcp listener.",
	)
	defer func() {
		_ = tcpListener.Close()
	}()
	tcpRaddr := tcpListener.Addr().String()

	unixNetwork := "unix"
	unixListener, e := net.Listen(unixNetwork, "")
	require.NoError(
		t,
		e,
		"DialNativeSyslog test expects no error when creating the"+
			" default unix listener.",
	)
	defer func() {
		_ = unixListener.Close()
	}()
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
			expectedError:     !syslogNewWorks,
			expectedSyslogger: syslogNewWorks,
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
			assert.Errorf(
				t,
				actualError,
				"DialNativeSyslog test expects an error for:"+
					" %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"DialNativeSyslog test expects no error for:"+
					" %s",
				explanation,
			)
		}

		if test.expectedSyslogger {
			assert.NotNil(
				t,
				actualSyslogger,
				"DialNativeSyslog test expects non-nil"+
					" syslogger for: %s",
				explanation,
			)
		} else {
			assert.Nil(
				t,
				actualSyslogger,
				"DialNativeSyslog test expects nil syslogger"+
					" for: %s",
				explanation,
			)
		}

		if actualSyslogger != nil {
			e := actualSyslogger.Close()
			assert.NoError(
				t,
				e,
				"NewNativeSyslog test expects no error when"+
					" closing the tested NativeSyslog for:"+
					" %s",
				explanation,
			)
		}
	}
}

func TestNativeSyslog(t *testing.T) {
	type content struct {
		S string
		E error
	}

	comm := make(chan *content)
	cancel := make(chan interface{})
	defer close(cancel)

	network := "unix"
	go func(comm chan<- *content, cancel <-chan interface{}) {
		defer close(comm)

		l, e := net.Listen(network, "")
		if e != nil {
			comm <- &content{E: e}
			return
		}
		go func(cancel <-chan interface{}, l net.Listener) {
			<-cancel
			_ = l.Close()
		}(cancel, l)
		raddr := l.Addr().String()
		comm <- &content{S: raddr}

		for {
			c, e := l.Accept()

			// If the cancellation channel has been closed, the
			// presumably that is why we'd see an Accept error, and
			// attempting to write this error to the comm channel
			// would panic.
			select {
			case <-cancel:
				return
			default:
			}

			if e != nil {
				comm <- &content{
					E: fmt.Errorf(
						"NativeSyslog test expects no"+
							" error when accepting"+
							" a new connection: %s",
						e.Error(),
					),
				}
				continue
			}

			r := bufio.NewReader(c)
			rs, e := r.ReadString('\n')
			if e != nil {
				comm <- &content{
					E: fmt.Errorf(
						"NativeSyslog test expects no"+
							" error when reading"+
							" with a syslog"+
							" reader: %s",
						e,
					),
				}
				_ = c.Close()
				continue
			}

			e = c.Close()
			if e != nil {
				comm <- &content{
					E: fmt.Errorf(
						"NativeSyslog test expects no"+
							" error when closing a"+
							" syslog reader: %s",
						e,
					),
				}
				continue
			}

			comm <- &content{S: rs[:len(rs)-1]}
		}
	}(comm, cancel)

	listenerInfo := <-comm
	require.NoError(
		t,
		listenerInfo.E,
		"NativeSyslog test requires no error from creation of"+
			" syslog listener",
	)
	raddr := listenerInfo.S

	facility := pri.User
	ident := "NativeSyslogTest"
	n, e := DialNativeSyslog(
		network,
		raddr,
		facility,
		ident,
	)
	require.NoError(
		t,
		e,
		"NativeSyslog test requires no error from creation of"+
			" NativeSyslog under test",
	)

	type testCase struct {
		inputPriority pri.Priority
		inputMsg      interface{}
		expectedError bool
	}

	tests := map[string]testCase{
		"nil values": {
			inputPriority: pri.Priority(0x0),
			inputMsg:      nil,
			expectedError: true,
		},
		"full values": {
			inputPriority: pri.Priority(0xFF),
			inputMsg:      "full values msg",
			expectedError: true,
		},
		"emerg": {
			inputPriority: pri.Emerg,
			inputMsg:      "emerg msg",
			expectedError: false,
		},
		"alert": {
			inputPriority: pri.Alert,
			inputMsg:      "alert msg",
			expectedError: false,
		},
		"crit": {
			inputPriority: pri.Crit,
			inputMsg:      "crit msg",
			expectedError: false,
		},
		"err": {
			inputPriority: pri.Err,
			inputMsg:      "err msg",
			expectedError: false,
		},
		"warning": {
			inputPriority: pri.Warning,
			inputMsg:      "warning msg",
			expectedError: false,
		},
		"notice": {
			inputPriority: pri.Notice,
			inputMsg:      "notice msg",
			expectedError: false,
		},
		"info": {
			inputPriority: pri.Info,
			inputMsg:      "info msg",
			expectedError: false,
		},
		"debug": {
			inputPriority: pri.Debug,
			inputMsg:      "debug msg",
			expectedError: false,
		},
		"combined priority": {
			inputPriority: pri.Syslog | pri.Notice,
			inputMsg:      "combined priority msg",
			expectedError: true,
		},
	}

	for explanation, test := range tests {
		actualError := n.Syslog(test.inputPriority, test.inputMsg)

		if test.expectedError {
			assert.Errorf(
				t,
				actualError,
				"NativeSyslog test expects an error for: %s",
				explanation,
			)
		} else {
			assert.NoError(
				t,
				actualError,
				"NativeSyslog test expects no error for: %s",
				explanation,
			)
		}

		// This assumes that if the error is nil then the listener will
		// definitely send on comm. If we just assume that it is safe to
		// proceed with this section if we don't expect an error, then
		// an error occurring which leads to no communication happening
		// will cause this test to hang forever.
		if actualError == nil {
			listenerResponse := <-comm
			require.NotNil(
				t,
				listenerResponse,
				"NativeSyslog test expects non-nil listener"+
					" response for: %s",
				explanation,
			)
			assert.NoError(
				t,
				listenerResponse.E,
				"NativeSyslog test expects no listener error"+
					" for: %s",
				explanation,
			)

			var expectedRegex *regexp.Regexp
			if s, ok := test.inputMsg.(string); ok {
				expectedRegex = regexp.MustCompile(
					fmt.Sprintf(
						"^<%d>.*%s$",
						facility|test.inputPriority,
						regexp.QuoteMeta(s),
					),
				)
			} else {
				expectedRegex = regexp.MustCompile(
					fmt.Sprintf(
						"^<%d>.*$",
						facility|test.inputPriority,
					),
				)
			}

			actualString := listenerResponse.S

			assert.Regexp(
				t,
				expectedRegex,
				actualString,
				"NativeSyslog test expects the raw syslog"+
					" message to indicate the correct"+
					" priority argument (and, if"+
					" applicable, have the correct"+
					" message) for: %s",
				explanation,
			)
		}
	}
}

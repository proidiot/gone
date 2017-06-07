package syslogger

import (
	"log/syslog"
	"os"
)

var osHostname = os.Hostname

var syslogNew = syslog.New

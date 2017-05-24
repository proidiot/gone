//go:generate stringer -type Severity
package pri

// Log severity
// See RFC 3164 Section 4.1.1, RFC 5424 Section 6.2.1, and RFC 5427 Section 3.
type Severity byte

const (
	// Emerg represents an emergency; system is unusable.
	Emerg Severity = 0
	// Alert represents than an action must be taken immediately.
	Alert Severity = 1
	// Crit represents a critical condition.
	Crit Severity = 2
	// Err represents an error condition.
	Err Severity = 3
	// Warning represents a warning condition.
	Warning Severity = 4
	// Notice represents a normal but significant condition.
	Notice Severity = 5
	// Info represents an informational message.
	Info Severity = 6
	// Debug represents debug-level messages.
	Debug Severity = 7
)

func (s Severity) Masked() Severity {
	// Equivalent to LOG_EMERG | LOG_ALERT | ... | LOG_DEBUG
	mask := Severity(0x07)
	return s & mask
}

package syslog

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestNewMessageDefaults(t *testing.T) {
	m := NewMessage("hello")

	if m.Facility != LOG_LOCAL0 {
		t.Errorf("Facility = %v, want LOG_LOCAL0", m.Facility)
	}
	if m.Severity != LOG_INFO {
		t.Errorf("Severity = %v, want LOG_INFO", m.Severity)
	}
	if m.Format != RFC3164 {
		t.Errorf("Format = %v, want RFC3164", m.Format)
	}
	if m.AppName != "echostrike" {
		t.Errorf("AppName = %q, want %q", m.AppName, "echostrike")
	}
	if m.Message != "hello" {
		t.Errorf("Message = %q, want %q", m.Message, "hello")
	}
}

func TestMessageStringRFC3164(t *testing.T) {
	ts := time.Date(2026, time.February, 14, 10, 0, 1, 0, time.UTC)
	m := &Message{
		Facility:  LOG_AUTH,
		Severity:  LOG_WARNING,
		Timestamp: ts,
		Hostname:  "server",
		AppName:   "sshd",
		ProcID:    "1234",
		Message:   "Failed password for root",
		Format:    RFC3164,
	}

	got := m.String()
	want := fmt.Sprintf("<%d>%s server sshd[1234]: Failed password for root",
		int(LOG_AUTH)|int(LOG_WARNING), ts.Format(time.Stamp))

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMessageStringRFC3164WithoutProcID(t *testing.T) {
	ts := time.Now()
	m := &Message{
		Facility:  LOG_LOCAL0,
		Severity:  LOG_INFO,
		Timestamp: ts,
		Hostname:  "server",
		AppName:   "echostrike",
		Message:   "test",
		Format:    RFC3164,
	}

	got := m.String()
	want := fmt.Sprintf("<%d>%s server echostrike: test", int(LOG_LOCAL0)|int(LOG_INFO), ts.Format(time.Stamp))
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMessageStringRFC5424(t *testing.T) {
	ts := time.Date(2026, time.February, 14, 10, 0, 1, 0, time.UTC)
	m := &Message{
		Facility:  LOG_LOCAL0,
		Severity:  LOG_INFO,
		Timestamp: ts,
		Hostname:  "server",
		AppName:   "echostrike",
		Message:   "hello",
		Format:    RFC5424,
	}

	got := m.String()
	want := fmt.Sprintf("<%d>1 %s server echostrike - - - hello",
		int(LOG_LOCAL0)|int(LOG_INFO), ts.Format(time.RFC3339))

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMessageStringRFC5424BlankFieldsBecomeDash(t *testing.T) {
	m := &Message{
		Facility:  LOG_LOCAL0,
		Severity:  LOG_INFO,
		Timestamp: time.Now(),
		Message:   "hello",
		Format:    RFC5424,
		// Hostname, AppName, ProcID, MsgID all left blank.
	}

	got := m.String()
	re := regexp.MustCompile(`^<\d+>1 \S+ - - - - - hello$`)
	if !re.MatchString(got) {
		t.Errorf("String() = %q, does not match expected RFC5424 blank-field pattern", got)
	}
}

func TestParseFacilityValid(t *testing.T) {
	cases := map[string]Priority{
		"kern":     LOG_KERN,
		"auth":     LOG_AUTH,
		"local0":   LOG_LOCAL0,
		"local7":   LOG_LOCAL7,
		"authpriv": LOG_AUTHPRIV,
	}

	for input, want := range cases {
		got, err := ParseFacility(input)
		if err != nil {
			t.Errorf("ParseFacility(%q) unexpected error: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("ParseFacility(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestParseFacilityInvalid(t *testing.T) {
	if _, err := ParseFacility("not-a-facility"); err == nil {
		t.Error("ParseFacility(\"not-a-facility\") expected error, got nil")
	}
}

func TestParseSeverityValid(t *testing.T) {
	cases := map[string]Priority{
		"emerg":   LOG_EMERG,
		"err":     LOG_ERR,
		"info":    LOG_INFO,
		"debug":   LOG_DEBUG,
		"warning": LOG_WARNING,
	}

	for input, want := range cases {
		got, err := ParseSeverity(input)
		if err != nil {
			t.Errorf("ParseSeverity(%q) unexpected error: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("ParseSeverity(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestParseSeverityInvalid(t *testing.T) {
	if _, err := ParseSeverity("not-a-severity"); err == nil {
		t.Error("ParseSeverity(\"not-a-severity\") expected error, got nil")
	}
}

// TestPriorityEncoding pins down the facility/severity packing scheme
// (facility pre-shifted by 3 bits, OR'd with a 0-7 severity) that the rest
// of the package's PRI values depend on.
func TestPriorityEncoding(t *testing.T) {
	m := NewMessage("x")
	m.Facility = LOG_LOCAL0 // 16 << 3 = 128
	m.Severity = LOG_ERR    // 3

	got := m.Facility | m.Severity
	want := Priority(131)
	if got != want {
		t.Errorf("Facility|Severity = %d, want %d", got, want)
	}
}

package syslog

import (
	"fmt"
	"os"
	"time"
)

type Priority int

const (
	LOG_KERN     Priority = 0 << 3
	LOG_USER     Priority = 1 << 3
	LOG_MAIL     Priority = 2 << 3
	LOG_DAEMON   Priority = 3 << 3
	LOG_AUTH     Priority = 4 << 3
	LOG_SYSLOG   Priority = 5 << 3
	LOG_LPR      Priority = 6 << 3
	LOG_NEWS     Priority = 7 << 3
	LOG_UUCP     Priority = 8 << 3
	LOG_CRON     Priority = 9 << 3
	LOG_AUTHPRIV Priority = 10 << 3
	LOG_FTP      Priority = 11 << 3
	LOG_LOCAL0   Priority = 16 << 3
	LOG_LOCAL1   Priority = 17 << 3
	LOG_LOCAL2   Priority = 18 << 3
	LOG_LOCAL3   Priority = 19 << 3
	LOG_LOCAL4   Priority = 20 << 3
	LOG_LOCAL5   Priority = 21 << 3
	LOG_LOCAL6   Priority = 22 << 3
	LOG_LOCAL7   Priority = 23 << 3
)

const (
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

type Format int

const (
	RFC3164 Format = iota
	RFC5424
)

type Message struct {
	Facility  Priority
	Severity  Priority
	Timestamp time.Time
	Hostname  string
	AppName   string
	ProcID    string
	MsgID     string
	Message   string
	Format    Format
}

func NewMessage(msg string) *Message {
	h, _ := os.Hostname()
	return &Message{
		Facility:  LOG_LOCAL0,
		Severity:  LOG_INFO,
		Timestamp: time.Now(),
		Hostname:  h,
		AppName:   "echostrike",
		Message:   msg,
		Format:    RFC3164, // Default to classic
	}
}

func (m *Message) String() string {
	p := m.Facility | m.Severity

	if m.Format == RFC5424 {
		// <PRIVAL>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID STRUCT-DATA MSG
		ts := m.Timestamp.Format(time.RFC3339)
		hostname := m.Hostname
		if hostname == "" {
			hostname = "-"
		}
		appName := m.AppName
		if appName == "" {
			appName = "-"
		}
		procID := m.ProcID
		if procID == "" {
			procID = "-"
		}
		msgID := m.MsgID
		if msgID == "" {
			msgID = "-"
		}

		return fmt.Sprintf("<%d>1 %s %s %s %s %s - %s",
			p, ts, hostname, appName, procID, msgID, m.Message)
	}

	ts := m.Timestamp.Format(time.Stamp)
	tag := m.AppName
	if m.ProcID != "" {
		tag = fmt.Sprintf("%s[%s]", m.AppName, m.ProcID)
	}
	return fmt.Sprintf("<%d>%s %s %s: %s", p, ts, m.Hostname, tag, m.Message)
}

func ParseFacility(f string) (Priority, error) {
	switch f {
	case "kern":
		return LOG_KERN, nil
	case "user":
		return LOG_USER, nil
	case "mail":
		return LOG_MAIL, nil
	case "daemon":
		return LOG_DAEMON, nil
	case "auth":
		return LOG_AUTH, nil
	case "syslog":
		return LOG_SYSLOG, nil
	case "lpr":
		return LOG_LPR, nil
	case "news":
		return LOG_NEWS, nil
	case "uucp":
		return LOG_UUCP, nil
	case "cron":
		return LOG_CRON, nil
	case "authpriv":
		return LOG_AUTHPRIV, nil
	case "ftp":
		return LOG_FTP, nil
	case "local0":
		return LOG_LOCAL0, nil
	case "local1":
		return LOG_LOCAL1, nil
	case "local2":
		return LOG_LOCAL2, nil
	case "local3":
		return LOG_LOCAL3, nil
	case "local4":
		return LOG_LOCAL4, nil
	case "local5":
		return LOG_LOCAL5, nil
	case "local6":
		return LOG_LOCAL6, nil
	case "local7":
		return LOG_LOCAL7, nil
	default:
		return 0, fmt.Errorf("unknown facility: %s", f)
	}
}

func ParseSeverity(s string) (Priority, error) {
	switch s {
	case "emerg":
		return LOG_EMERG, nil
	case "alert":
		return LOG_ALERT, nil
	case "crit":
		return LOG_CRIT, nil
	case "err":
		return LOG_ERR, nil
	case "warning":
		return LOG_WARNING, nil
	case "notice":
		return LOG_NOTICE, nil
	case "info":
		return LOG_INFO, nil
	case "debug":
		return LOG_DEBUG, nil
	default:
		return 0, fmt.Errorf("unknown severity: %s", s)
	}
}

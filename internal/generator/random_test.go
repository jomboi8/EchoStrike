package generator

import (
	"regexp"
	"slices"
	"testing"
)

func TestRandomIPFormat(t *testing.T) {
	re := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	for range 100 {
		ip := RandomIP()
		if !re.MatchString(ip) {
			t.Fatalf("RandomIP() = %q, does not look like an IPv4 address", ip)
		}
	}
}

func TestRandomUserIsFromList(t *testing.T) {
	for range 100 {
		u := RandomUser()
		if !slices.Contains(userNames, u) {
			t.Fatalf("RandomUser() = %q, not in known user list %v", u, userNames)
		}
	}
}

func TestRandomPortInRange(t *testing.T) {
	for range 100 {
		p := RandomPort()
		if p < 1024 || p >= 65535 {
			t.Fatalf("RandomPort() = %d, want value in [1024, 65535)", p)
		}
	}
}

func TestRandomStatusCodeIsKnown(t *testing.T) {
	known := []int{200, 201, 301, 302, 400, 401, 403, 404, 500, 502, 503}
	for range 100 {
		code := RandomStatusCode()
		if !slices.Contains(known, code) {
			t.Fatalf("RandomStatusCode() = %d, not in known status list %v", code, known)
		}
	}
}

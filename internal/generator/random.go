package generator

import (
	"fmt"
	"math/rand"
)

var (
	userNames = []string{"admin", "root", "user", "guest", "support", "sysadmin", "deploy"}
	actions   = []string{"login", "logout", "failed", "accepted", "error", "timeout"}
	paths     = []string{"/index.html", "/api/v1/login", "/admin", "/dashboard", "/config", "/etc/passwd"}
)

// RandomIP generates a random IPv4 address
func RandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

// RandomUser returns a random username
func RandomUser() string {
	return userNames[rand.Intn(len(userNames))]
}

// RandomAction returns a random action string
func RandomAction() string {
	return actions[rand.Intn(len(actions))]
}

// RandomPath returns a random file/url path
func RandomPath() string {
	return paths[rand.Intn(len(paths))]
}

// RandomStatusCode returns a random HTTP status code
func RandomStatusCode() int {
	codes := []int{200, 201, 301, 302, 400, 401, 403, 404, 500, 502, 503}
	return codes[rand.Intn(len(codes))]
}

// RandomPort returns a random port number
func RandomPort() int {
	return rand.Intn(65535-1024) + 1024
}

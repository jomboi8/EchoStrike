package generator

import (
	"bytes"
	"text/template"
	"time"
)

// TemplateData holds the dynamic data for log generation
type TemplateData struct {
	Timestamp string
	IP        string
	User      string
	Action    string
	Path      string
	Status    int
	Port      int
}

// Generator manages log templates and generation
type Generator struct {
	templates map[string]*template.Template
}

func New() *Generator {
	g := &Generator{
		templates: make(map[string]*template.Template),
	}
	g.registerDefaults()
	return g
}

func (g *Generator) registerDefaults() {
	// SSH Failed Login
	g.Register("ssh-failed", "Failed password for {{.User}} from {{.IP}} port {{.Port}} ssh2")

	// SSH Accepted
	g.Register("ssh-accepted", "Accepted publickey for {{.User}} from {{.IP}} port {{.Port}} ssh2")

	// Nginx Access
	g.Register("nginx-access", `{{.IP}} - {{.User}} [{{.Timestamp}}] "GET {{.Path}} HTTP/1.1" {{.Status}} 1024 "-" "Mozilla/5.0"`)

	// Firewall Drop
	g.Register("firewall-drop", "SRC={{.IP}} DST=192.168.1.1 PROTO=TCP DPT={{.Port}} ACTION=DROP")
}

func (g *Generator) Register(name, tmplString string) {
	tmpl, err := template.New(name).Parse(tmplString)
	if err == nil {
		g.templates[name] = tmpl
	}
}

func (g *Generator) Generate(templateName string) (string, error) {
	tmpl, ok := g.templates[templateName]
	if !ok {
		return "", getTemplateError(templateName)
	}

	data := TemplateData{
		Timestamp: time.Now().Format(time.RFC3339),
		IP:        RandomIP(),
		User:      RandomUser(),
		Action:    RandomAction(),
		Path:      RandomPath(),
		Status:    RandomStatusCode(),
		Port:      RandomPort(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getTemplateError(name string) error {
	return &TemplateNotFoundError{Name: name}
}

type TemplateNotFoundError struct {
	Name string
}

func (e *TemplateNotFoundError) Error() string {
	return "template not found: " + e.Name
}

func (g *Generator) ListTemplates() []string {
	keys := make([]string, 0, len(g.templates))
	for k := range g.templates {
		keys = append(keys, k)
	}
	return keys
}

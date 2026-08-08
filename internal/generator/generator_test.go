package generator

import (
	"errors"
	"regexp"
	"slices"
	"testing"
)

func TestGenerateKnownTemplates(t *testing.T) {
	g := New()

	patterns := map[string]*regexp.Regexp{
		"ssh-failed":    regexp.MustCompile(`^Failed password for \S+ from \d+\.\d+\.\d+\.\d+ port \d+ ssh2$`),
		"ssh-accepted":  regexp.MustCompile(`^Accepted publickey for \S+ from \d+\.\d+\.\d+\.\d+ port \d+ ssh2$`),
		"firewall-drop": regexp.MustCompile(`^SRC=\d+\.\d+\.\d+\.\d+ DST=192\.168\.1\.1 PROTO=TCP DPT=\d+ ACTION=DROP$`),
	}

	for name, pattern := range patterns {
		got, err := g.Generate(name)
		if err != nil {
			t.Errorf("Generate(%q) unexpected error: %v", name, err)
			continue
		}
		if !pattern.MatchString(got) {
			t.Errorf("Generate(%q) = %q, does not match expected pattern %s", name, got, pattern)
		}
	}
}

func TestGenerateUnknownTemplateReturnsTypedError(t *testing.T) {
	g := New()

	_, err := g.Generate("does-not-exist")
	if err == nil {
		t.Fatal("Generate() with unknown template expected error, got nil")
	}

	var notFound *TemplateNotFoundError
	if !errors.As(err, &notFound) {
		t.Errorf("Generate() error = %T, want *TemplateNotFoundError", err)
	}
}

func TestRegisterCustomTemplate(t *testing.T) {
	g := New()
	g.Register("custom", "user={{.User}} action={{.Action}}")

	got, err := g.Generate("custom")
	if err != nil {
		t.Fatalf("Generate(\"custom\") unexpected error: %v", err)
	}

	want := regexp.MustCompile(`^user=\S+ action=\S+$`)
	if !want.MatchString(got) {
		t.Errorf("Generate(\"custom\") = %q, does not match expected pattern", got)
	}
}

func TestRegisterInvalidTemplateIsIgnored(t *testing.T) {
	g := New()
	before := len(g.ListTemplates())

	g.Register("broken", "{{.Unclosed")

	if len(g.ListTemplates()) != before {
		t.Error("Register() with invalid template text should not register anything")
	}
}

func TestListTemplatesIncludesDefaults(t *testing.T) {
	g := New()
	names := g.ListTemplates()

	want := []string{"ssh-failed", "ssh-accepted", "nginx-access", "firewall-drop"}
	for _, w := range want {
		if !slices.Contains(names, w) {
			t.Errorf("ListTemplates() = %v, missing default template %q", names, w)
		}
	}
}

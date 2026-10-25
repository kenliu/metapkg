package packages

import (
	"strings"
	"testing"
)

func TestParseOSRelease(t *testing.T) {
	input := `
PRETTY_NAME="Ubuntu 24.04.5 LTS"
NAME="Ubuntu"
# a comment
ID=ubuntu
ID_LIKE=debian
VERSION_CODENAME='noble'
MALFORMED LINE
`
	got := parseOSRelease(strings.NewReader(input))

	expected := map[string]string{
		"PRETTY_NAME":      "Ubuntu 24.04.5 LTS",
		"NAME":             "Ubuntu",
		"ID":               "ubuntu",
		"ID_LIKE":          "debian",
		"VERSION_CODENAME": "noble",
	}

	if len(got) != len(expected) {
		t.Errorf("expected %d entries, got %d: %v", len(expected), len(got), got)
	}
	for key, want := range expected {
		if got[key] != want {
			t.Errorf("expected %s=%q, got %q", key, want, got[key])
		}
	}
}

func TestDefaultPackageManagerForOSRelease(t *testing.T) {
	tests := []struct {
		name      string
		osRelease map[string]string
		expected  string
	}{
		{"ubuntu", map[string]string{"ID": "ubuntu", "ID_LIKE": "debian"}, "apt"},
		{"debian", map[string]string{"ID": "debian"}, "apt"},
		{"linux mint via ID_LIKE", map[string]string{"ID": "linuxmint", "ID_LIKE": "ubuntu debian"}, "apt"},
		{"pop os via ID_LIKE", map[string]string{"ID": "pop", "ID_LIKE": "ubuntu debian"}, "apt"},
		{"fedora", map[string]string{"ID": "fedora"}, "dnf"},
		{"rocky via ID_LIKE", map[string]string{"ID": "rocky", "ID_LIKE": "rhel centos fedora"}, "dnf"},
		{"unknown distro falls back to default", map[string]string{"ID": "arch"}, DEFAULT_PACKAGE_MANAGER},
		{"empty os-release falls back to default", map[string]string{}, DEFAULT_PACKAGE_MANAGER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultPackageManagerForOSRelease(tt.osRelease); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

package script

import (
	"reflect"
	"testing"

	"github.com/kenliu/metapkg/internal/packages"
)

func TestInstallDescription(t *testing.T) {
	tests := []struct {
		name      string
		scriptdef packages.Scriptdef
		expected  []string
	}{
		{
			name: "one command per line",
			scriptdef: packages.Scriptdef{Name: "vivaldi", Commands: []string{
				"sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo",
				"sudo dnf install -y vivaldi-stable",
			}},
			expected: []string{
				"sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo",
				"sudo dnf install -y vivaldi-stable",
			},
		},
		{
			name: "multi-line command is split into lines",
			scriptdef: packages.Scriptdef{Name: "gcloud", Commands: []string{
				"sudo tee /etc/yum.repos.d/google-cloud-sdk.repo << EOM\n[google-cloud-cli]\nEOM",
				"sudo dnf install google-cloud-cli",
			}},
			expected: []string{
				"sudo tee /etc/yum.repos.d/google-cloud-sdk.repo << EOM",
				"[google-cloud-cli]",
				"EOM",
				"sudo dnf install google-cloud-cli",
			},
		},
		{
			name:      "no commands",
			scriptdef: packages.Scriptdef{Name: "empty"},
			expected:  []string{"(no commands found for scriptdef empty)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScriptPackageState(tt.scriptdef)
			if got := s.InstallDescription(tt.scriptdef.Name); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("InstallDescription(%q) = %q, expected %q", tt.scriptdef.Name, got, tt.expected)
			}
		})
	}
}

func TestScriptContent(t *testing.T) {
	s := NewScriptPackageState(packages.Scriptdef{Name: "x", Commands: []string{"echo one", "echo two"}})
	expected := "#!/bin/bash\n\necho one\necho two"
	if got := s.scriptContent(); got != expected {
		t.Errorf("scriptContent() = %q, expected %q", got, expected)
	}
}

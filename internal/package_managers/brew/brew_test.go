package brew

import (
	"strings"
	"testing"
)

func TestInstallCommand(t *testing.T) {
	tests := []struct {
		name     string
		state    BrewPackageState
		pkg      string
		expected string
	}{
		{"kind not yet known", BrewPackageState{}, "bat", "brew install bat"},
		{"formula", BrewPackageState{kind: kindFormula}, "bat", "brew install --formula bat"},
		{"cask", BrewPackageState{kind: kindCask}, "iterm2", "brew install --cask iterm2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strings.Join(tt.state.InstallCommand(tt.pkg), " ")
			if got != tt.expected {
				t.Errorf("InstallCommand(%q) = %q, expected %q", tt.pkg, got, tt.expected)
			}
		})
	}
}

func TestKindFromInfo(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		expected  packageKind
		expectErr bool
	}{
		{"formula", `{"formulae":[{"name":"bat"}],"casks":[]}`, kindFormula, false},
		{"cask", `{"formulae":[],"casks":[{"token":"iterm2"}]}`, kindCask, false},
		{"both prefers formula", `{"formulae":[{"name":"docker"}],"casks":[{"token":"docker"}]}`, kindFormula, false},
		{"neither", `{"formulae":[],"casks":[]}`, kindUnknown, true},
		{"invalid json", `not json`, kindUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kindFromInfo("pkg", []byte(tt.json))
			if (err != nil) != tt.expectErr {
				t.Fatalf("kindFromInfo error = %v, expected error: %v", err, tt.expectErr)
			}
			if got != tt.expected {
				t.Errorf("kindFromInfo = %v, expected %v", got, tt.expected)
			}
		})
	}
}

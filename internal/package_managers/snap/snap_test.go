package snap

import "testing"

func TestIsNotInstalledMessage(t *testing.T) {
	tests := []struct {
		name     string
		stderr   string
		expected bool
	}{
		{"not installed", "error: no matching snaps installed\n", true},
		{"snapd not running", "error: cannot communicate with server: Get \"http://localhost/v2/snaps\": dial unix /run/snapd.socket: connect: no such file or directory\n", false},
		{"empty output", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNotInstalledMessage(tt.stderr); got != tt.expected {
				t.Errorf("isNotInstalledMessage(%q) = %v, expected %v", tt.stderr, got, tt.expected)
			}
		})
	}
}

func TestIsOutdatedOutput(t *testing.T) {
	refreshOutput := `Name     Version        Rev    Size   Publisher   Notes
core22   20240111       1122   77MB   canonical✓  base
firefox  122.0-2        3836   262MB  mozilla✓    -
`
	upToDateOutput := `All snaps up to date.
`

	tests := []struct {
		name     string
		output   string
		pkg      string
		expected bool
	}{
		{"update available", refreshOutput, "firefox", true},
		{"update available for base snap", refreshOutput, "core22", true},
		{"different snap outdated", refreshOutput, "chromium", false},
		{"partial name does not match", refreshOutput, "fire", false},
		{"header is not a snap", refreshOutput, "Name", false},
		{"all snaps up to date", upToDateOutput, "firefox", false},
		{"all snaps up to date message is not a snap", upToDateOutput, "All", false},
		{"empty output", "", "firefox", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOutdatedOutput(tt.output, tt.pkg); got != tt.expected {
				t.Errorf("isOutdatedOutput(..., %q) = %v, expected %v", tt.pkg, got, tt.expected)
			}
		})
	}
}

func TestRequiresClassic(t *testing.T) {
	classicOutput := `error: This revision of snap "code" was published using classic confinement and thus may perform
       arbitrary system changes outside of the security sandbox that snaps are usually confined to,
       which may put your system at risk.

       If you understand and want to proceed repeat the command including --classic.
`
	notFoundOutput := `error: snap "no-such-snap" not found
`

	tests := []struct {
		name     string
		stderr   string
		expected bool
	}{
		{"classic confinement required", classicOutput, true},
		{"snap not found", notFoundOutput, false},
		{"empty output", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requiresClassic(tt.stderr); got != tt.expected {
				t.Errorf("requiresClassic(%q) = %v, expected %v", tt.stderr, got, tt.expected)
			}
		})
	}
}

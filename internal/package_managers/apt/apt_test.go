package apt

import (
	"strings"
	"testing"
)

func TestIsInstalledStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"installed", "install ok installed", true},
		{"installed with trailing newline", "install ok installed\n", true},
		{"removed but config files remain", "deinstall ok config-files", false},
		{"unpacked but not configured", "install ok unpacked", false},
		{"half-installed", "install ok half-installed", false},
		{"not installed", "unknown ok not-installed", false},
		{"empty output", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInstalledStatus(tt.status); got != tt.expected {
				t.Errorf("isInstalledStatus(%q) = %v, expected %v", tt.status, got, tt.expected)
			}
		})
	}
}

func TestInstallCommand(t *testing.T) {
	got := strings.Join((&AptPackageState{}).InstallCommand("curl"), " ")
	expected := "sudo apt-get install -y curl"
	if got != expected {
		t.Errorf("InstallCommand(\"curl\") = %q, expected %q", got, expected)
	}
}

func TestIsOutdatedOutput(t *testing.T) {
	upgradeOutput := `NOTE: This is only a simulation!
      apt-get needs root privileges for real execution.
      Keep also in mind that locking is deactivated,
      so don't depend on the relevance to the real current situation!
Reading package lists...
Building dependency tree...
Reading state information...
The following packages will be upgraded:
  curl
1 upgraded, 0 newly installed, 0 to remove and 12 not upgraded.
Inst curl [8.5.0-2ubuntu10.5] (8.5.0-2ubuntu10.6 Ubuntu:24.04/noble-updates [amd64])
Conf curl (8.5.0-2ubuntu10.6 Ubuntu:24.04/noble-updates [amd64])
`
	upToDateOutput := `NOTE: This is only a simulation!
Reading package lists...
Building dependency tree...
Reading state information...
curl is already the newest version (8.5.0-2ubuntu10.6).
0 upgraded, 0 newly installed, 0 to remove and 12 not upgraded.
`

	tests := []struct {
		name     string
		output   string
		pkg      string
		expected bool
	}{
		{"upgrade available", upgradeOutput, "curl", true},
		{"already newest version", upToDateOutput, "curl", false},
		{"different package upgraded", upgradeOutput, "wget", false},
		{"empty output", "", "curl", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOutdatedOutput(tt.output, tt.pkg); got != tt.expected {
				t.Errorf("isOutdatedOutput(..., %q) = %v, expected %v", tt.pkg, got, tt.expected)
			}
		})
	}
}

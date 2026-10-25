package engine

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/kenliu/metapkg/internal/package_managers"
	"github.com/kenliu/metapkg/internal/packages"
)

// fakePackageState is a PackageState that records whether Install was called
// and does not touch the system. It does not implement InstallCommander or
// InstallDescriber.
type fakePackageState struct {
	installed      bool
	isInstalledErr error
	installCalled  bool
}

func (f *fakePackageState) IsInstalled(name string, arguments []string) (bool, error) {
	return f.installed, f.isInstalledErr
}

func (f *fakePackageState) IsOutdated(name string, arguments []string) (bool, error) {
	return false, nil
}

func (f *fakePackageState) Install(name string) error {
	f.installCalled = true
	return nil
}

// fakeCommander is a fakePackageState that also describes a single install
// command, like apt, dnf and brew.
type fakeCommander struct {
	fakePackageState
	command []string
}

func (f *fakeCommander) InstallCommand(name string) []string {
	return f.command
}

// fakeRetryingCommander is a fakeCommander that also describes a conditional
// retry, like snap and flatpak.
type fakeRetryingCommander struct {
	fakeCommander
	details []string
}

func (f *fakeRetryingCommander) InstallDescription(name string) []string {
	return f.details
}

// fakeScripter is a fakePackageState that only describes the commands it
// would run, like script.
type fakeScripter struct {
	fakePackageState
	details []string
}

func (f *fakeScripter) InstallDescription(name string) []string {
	return f.details
}

// useFakes replaces newPackageState for the duration of the test with one
// that returns the fake registered for each package name. Packages with no
// registered fake fall through to the real constructor, so unsupported
// package managers still produce the real error.
func useFakes(t *testing.T, fakes map[string]package_managers.PackageState) {
	t.Helper()
	original := newPackageState
	newPackageState = func(m *packages.Metapackage, pkg packages.Package) (package_managers.PackageState, error) {
		if state, ok := fakes[pkg.Name]; ok {
			return state, nil
		}
		return original(m, pkg)
	}
	t.Cleanup(func() { newPackageState = original })
}

func TestDryRunInstall(t *testing.T) {
	tests := []struct {
		name          string
		packages      []packages.Package
		fakes         map[string]package_managers.PackageState
		expectedLines []string
		expectedErr   string
	}{
		{
			name: "all installed",
			packages: []packages.Package{
				{Name: "curl", PackageManager: "apt"},
				{Name: "htop", PackageManager: "apt"},
			},
			fakes: map[string]package_managers.PackageState{
				"curl": &fakeCommander{fakePackageState: fakePackageState{installed: true}},
				"htop": &fakeCommander{fakePackageState: fakePackageState{installed: true}},
			},
			expectedLines: []string{
				"Dry run: no changes will be made",
				"[skip] curl (apt): already installed",
				"[skip] htop (apt): already installed",
				"Summary: 0 to install, 2 already installed",
				"Dry run: nothing was changed",
			},
		},
		{
			name: "none installed",
			packages: []packages.Package{
				{Name: "sl", PackageManager: "apt"},
				{Name: "bat", PackageManager: "brew"},
			},
			fakes: map[string]package_managers.PackageState{
				"sl":  &fakeCommander{command: []string{"sudo", "apt-get", "install", "-y", "sl"}},
				"bat": &fakeCommander{command: []string{"brew", "install", "bat"}},
			},
			expectedLines: []string{
				"[install] sl (apt): would run: sudo apt-get install -y sl",
				"[install] bat (brew): would run: brew install bat",
				"Summary: 2 to install, 0 already installed",
			},
		},
		{
			name: "mixed",
			packages: []packages.Package{
				{Name: "curl", PackageManager: "apt"},
				{Name: "starship", PackageManager: "snap"},
				{Name: "htop", PackageManager: "apt"},
			},
			fakes: map[string]package_managers.PackageState{
				"curl": &fakeCommander{fakePackageState: fakePackageState{installed: true}},
				"starship": &fakeRetryingCommander{
					fakeCommander: fakeCommander{command: []string{"sudo", "snap", "install", "starship"}},
					details:       []string{"may retry with: sudo snap install --classic starship (if snap reports the snap requires classic confinement)"},
				},
				"htop": &fakeCommander{fakePackageState: fakePackageState{installed: true}},
			},
			expectedLines: []string{
				"[skip] curl (apt): already installed",
				"[install] starship (snap): would run: sudo snap install starship",
				"    may retry with: sudo snap install --classic starship (if snap reports the snap requires classic confinement)",
				"[skip] htop (apt): already installed",
				"Summary: 1 to install, 2 already installed",
			},
		},
		{
			name: "script package describes its commands",
			packages: []packages.Package{
				{Name: "vivaldi", PackageManager: "script"},
			},
			fakes: map[string]package_managers.PackageState{
				"vivaldi": &fakeScripter{details: []string{
					"sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo",
					"sudo dnf install -y vivaldi-stable",
				}},
			},
			expectedLines: []string{
				"[install] vivaldi (script): would run:",
				"    sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo",
				"    sudo dnf install -y vivaldi-stable",
				"Summary: 1 to install, 0 already installed",
			},
		},
		{
			name: "package manager without a command description",
			packages: []packages.Package{
				{Name: "foo", PackageManager: "apt"},
			},
			fakes: map[string]package_managers.PackageState{
				"foo": &fakePackageState{},
			},
			expectedLines: []string{
				"[install] foo (apt): would install",
				"Summary: 1 to install, 0 already installed",
			},
		},
		{
			name: "IsInstalled error propagates",
			packages: []packages.Package{
				{Name: "curl", PackageManager: "apt"},
				{Name: "starship", PackageManager: "snap"},
			},
			fakes: map[string]package_managers.PackageState{
				"curl":     &fakeCommander{fakePackageState: fakePackageState{installed: true}},
				"starship": &fakeCommander{fakePackageState: fakePackageState{isInstalledErr: errors.New("snap: executable not found")}},
			},
			expectedErr: "snap: executable not found",
		},
		{
			name: "unsupported package manager errors",
			packages: []packages.Package{
				{Name: "curl", PackageManager: "apt"},
				{Name: "vim", PackageManager: "pacman"},
			},
			fakes: map[string]package_managers.PackageState{
				"curl": &fakeCommander{fakePackageState: fakePackageState{installed: true}},
			},
			expectedErr: "unsupported package manager: pacman",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useFakes(t, tt.fakes)
			m := &packages.Metapackage{Packages: tt.packages, Scriptdefs: map[string]packages.Scriptdef{}}

			var out strings.Builder
			err := DryRunInstall(m, &out)

			if tt.expectedErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("DryRunInstall() error = %v, expected error containing %q", err, tt.expectedErr)
				}
				if out.Len() != 0 {
					t.Errorf("DryRunInstall() printed a plan despite the error:\n%s", out.String())
				}
			} else {
				if err != nil {
					t.Fatalf("DryRunInstall() unexpected error: %v", err)
				}
				got := out.String()
				for _, line := range tt.expectedLines {
					if !strings.Contains(got, line+"\n") {
						t.Errorf("DryRunInstall() output missing line %q\noutput:\n%s", line, got)
					}
				}
			}

			// Install must never be called in a dry run, whatever the outcome
			for name, state := range tt.fakes {
				if installCalled(state) {
					t.Errorf("DryRunInstall() called Install for %s", name)
				}
			}
		})
	}
}

// installCalled reports whether Install was called on one of the fakes.
func installCalled(state package_managers.PackageState) bool {
	switch s := state.(type) {
	case *fakePackageState:
		return s.installCalled
	case *fakeCommander:
		return s.installCalled
	case *fakeRetryingCommander:
		return s.installCalled
	case *fakeScripter:
		return s.installCalled
	}
	return false
}

func TestDryRunInstallPlanIsInFileOrder(t *testing.T) {
	useFakes(t, map[string]package_managers.PackageState{
		"c": &fakeCommander{fakePackageState: fakePackageState{installed: true}},
		"a": &fakeCommander{command: []string{"brew", "install", "a"}},
		"b": &fakeCommander{fakePackageState: fakePackageState{installed: true}},
	})
	m := &packages.Metapackage{Packages: []packages.Package{
		{Name: "c", PackageManager: "brew"},
		{Name: "a", PackageManager: "brew"},
		{Name: "b", PackageManager: "brew"},
	}}

	var out strings.Builder
	if err := DryRunInstall(m, &out); err != nil {
		t.Fatalf("DryRunInstall() unexpected error: %v", err)
	}

	expected := `
Dry run: no changes will be made
================================
[skip] c (brew): already installed
[install] a (brew): would run: brew install a
[skip] b (brew): already installed

Summary: 1 to install, 2 already installed
Dry run: nothing was changed
`
	if out.String() != expected {
		t.Errorf("DryRunInstall() output =\n%s\nexpected:\n%s", out.String(), expected)
	}
}

// TestDryRunInstallRealScriptPackage exercises the real script package
// manager through the engine. "which" is read-only and the name is chosen so
// it is never on the PATH.
func TestDryRunInstallRealScriptPackage(t *testing.T) {
	if _, err := exec.LookPath("which"); err != nil {
		t.Skip("which is not available")
	}

	name := "metapkg-dry-run-test-not-on-path"
	m := &packages.Metapackage{
		Packages: []packages.Package{{Name: name, PackageManager: "script"}},
		Scriptdefs: map[string]packages.Scriptdef{
			name: {Name: name, Commands: []string{"echo one", "echo two"}},
		},
	}

	var out strings.Builder
	if err := DryRunInstall(m, &out); err != nil {
		t.Fatalf("DryRunInstall() unexpected error: %v", err)
	}

	expectedLines := []string{
		"[install] " + name + " (script): would run:",
		"    echo one",
		"    echo two",
		"Summary: 1 to install, 0 already installed",
	}
	for _, line := range expectedLines {
		if !strings.Contains(out.String(), line+"\n") {
			t.Errorf("DryRunInstall() output missing line %q\noutput:\n%s", line, out.String())
		}
	}
}

func TestInstallPackagesCallsInstall(t *testing.T) {
	installed := &fakeCommander{fakePackageState: fakePackageState{installed: true}}
	missing := &fakeCommander{command: []string{"brew", "install", "missing"}}
	useFakes(t, map[string]package_managers.PackageState{
		"installed": installed,
		"missing":   missing,
	})
	m := &packages.Metapackage{Packages: []packages.Package{
		{Name: "installed", PackageManager: "brew"},
		{Name: "missing", PackageManager: "brew"},
	}}

	if err := InstallPackages(m); err != nil {
		t.Fatalf("InstallPackages() unexpected error: %v", err)
	}
	if installed.installCalled {
		t.Errorf("InstallPackages() called Install for an already installed package")
	}
	if !missing.installCalled {
		t.Errorf("InstallPackages() did not call Install for a missing package")
	}
}

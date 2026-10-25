package brew

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// packageKind says whether a name refers to a Homebrew formula or a cask.
type packageKind int

const (
	kindUnknown packageKind = iota
	kindFormula
	kindCask
)

// flag returns the brew option that selects this kind of package.
func (k packageKind) flag() string {
	if k == kindCask {
		return "--cask"
	}
	return "--formula"
}

// BrewPackageState implements the PackageState interface for the Homebrew package manager.
// A name can be a formula or a cask. IsInstalled works out which one it is and
// remembers the answer, so Install and InstallCommand can pass the right option.
type BrewPackageState struct {
	kind packageKind
}

func (b *BrewPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	for _, kind := range []packageKind{kindFormula, kindCask} {
		installed, err := listed(name, kind)
		if err != nil {
			return false, err
		}
		if installed {
			b.kind = kind
			return true, nil
		}
	}

	// Not installed, so find out whether it is a formula or a cask for Install.
	kind, err := lookupKind(name)
	if err != nil {
		return false, err
	}
	b.kind = kind
	return false, nil
}

func (b *BrewPackageState) IsOutdated(name string, arguments []string) (bool, error) {
	installed, err := b.IsInstalled(name, arguments)
	if err != nil {
		return false, err
	}
	if !installed {
		return false, nil
	}

	fmt.Printf("Executing command: brew outdated %s %s\n", b.kind.flag(), name)

	cmd := exec.Command("brew", "outdated", b.kind.flag(), name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return false, fmt.Errorf("error checking if package is outdated: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	// If the output is not empty, the package is outdated
	return strings.TrimSpace(stdoutStr) != "", nil
}

// InstallCommand returns the command Install runs, without running it.
// Once IsInstalled has run it names the kind explicitly (--formula or --cask).
func (b *BrewPackageState) InstallCommand(name string) []string {
	if b.kind == kindUnknown {
		return []string{"brew", "install", name}
	}
	return []string{"brew", "install", b.kind.flag(), name}
}

func (b *BrewPackageState) Install(name string) error {
	args := b.InstallCommand(name)
	fmt.Printf("Executing command: %s\n", strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return fmt.Errorf("error installing package: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}
	return nil
}

// listed reports whether name is installed as the given kind of package.
func listed(name string, kind packageKind) (bool, error) {
	fmt.Printf("Executing command: brew list %s %s\n", kind.flag(), name)

	cmd := exec.Command("brew", "list", kind.flag(), name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		// brew will return with error code 1 if the package is not installed
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("error checking if package is installed: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	return true, nil
}

// lookupKind asks Homebrew whether name is a formula or a cask.
func lookupKind(name string) (packageKind, error) {
	fmt.Printf("Executing command: brew info --json=v2 %s\n", name)

	cmd := exec.Command("brew", "info", "--json=v2", name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return kindUnknown, fmt.Errorf("%q is not a Homebrew formula or cask: %w\nstderr: %s", name, err, stderr.String())
	}

	return kindFromInfo(name, []byte(stdout.String()))
}

// kindFromInfo reads the output of `brew info --json=v2` and decides whether
// the package is a formula or a cask. A name that is both is an error, because
// metapkg cannot tell which one the file means.
func kindFromInfo(name string, data []byte) (packageKind, error) {
	var info struct {
		Formulae []json.RawMessage `json:"formulae"`
		Casks    []json.RawMessage `json:"casks"`
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return kindUnknown, fmt.Errorf("error reading brew info output for %q: %w", name, err)
	}

	isFormula, isCask := len(info.Formulae) > 0, len(info.Casks) > 0
	switch {
	case isFormula && isCask:
		return kindUnknown, fmt.Errorf("%q is ambiguous: it is both a Homebrew formula and a cask", name)
	case isFormula:
		return kindFormula, nil
	case isCask:
		return kindCask, nil
	}
	return kindUnknown, fmt.Errorf("%q is not a Homebrew formula or cask", name)
}

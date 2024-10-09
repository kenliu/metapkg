package brew

import (
	"fmt"
	"os/exec"
	"strings"
)

// BrewPackageState implements the PackageState interface for the Homebrew package manager
type BrewPackageState struct{}

func (b *BrewPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	// Execute the brew command to check if the specific package is installed
	fmt.Printf("Executing command: brew list --formula %s\n", name)

	cmd := exec.Command("brew", "list", "--formula", name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Capture stdout and stderr as strings
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		// brew will return with error code 1 if the package is not installed
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("error checking if package is installed: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	// If we get here, the package is installed
	return true, nil
}

func (b *BrewPackageState) IsOutdated(name string, arguments []string) (bool, error) {
	fmt.Printf("Executing command: brew outdated --formula %s\n", name)

	cmd := exec.Command("brew", "outdated", "--formula", name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return false, fmt.Errorf("error checking if package is outdated: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	// If the output is not empty, the package is outdated
	return strings.TrimSpace(stdoutStr) != "", nil
}

func (b *BrewPackageState) Install(name string) error {
	fmt.Printf("Executing command: brew install %s\n", name)

	cmd := exec.Command("brew", "install", name)
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

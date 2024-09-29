package package_managers

import (
	"fmt"
	"os/exec"
	"strings"
)

// DnfPackageState implements the PackageState interface for the DNF package manager
type DnfPackageState struct {
}

func (d *DnfPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	// execute the dnf command to check if the specific package is installed
	// return the output of the command
	// if the package is installed, return true
	// if the package is not installed, return false
	// if there is an error, return the error

	// log the command to be executed
	fmt.Printf("Executing command: dnf list installed %s\n", name)

	cmd := exec.Command("dnf", append([]string{"list", "installed", name})...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Capture stdout and stderr as strings
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		// dnf will return with error code 1 if the package is not installed
		// https://dnf.readthedocs.io/en/latest/command_ref.html#description
		packageNotInstalledMessage := "Error: No matching Packages to list"
		if strings.Contains(stderrStr, packageNotInstalledMessage) {
			return false, nil
		} else {
			return false, fmt.Errorf("error checking if package is installed: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
		}
	}

	// If the package is installed, the output will contain more than just the header line
	lines := strings.Split(strings.TrimSpace(stdoutStr), "\n")
	// find the line that contains the package name
	for _, line := range lines {
		if strings.Contains(line, name) {
			return true, nil
		}
	}
	return false, nil
}

func (d *DnfPackageState) IsOutdated(name string, arguments []string) (bool, error) {
	return false, nil
}

func (d *DnfPackageState) Install(name string) error {
	// execute the dnf command to install the specific package
	// return the output of the command
	// if there is an error, return the error

	// log the command to be executed
	fmt.Printf("Executing command: dnf install %s\n", name)

	cmd := exec.Command("sudo", append([]string{"dnf", "install", "-y", name})...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Capture stdout and stderr as strings
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return fmt.Errorf("error installing package: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}
	return nil
}

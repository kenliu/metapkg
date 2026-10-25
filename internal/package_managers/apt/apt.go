package apt

import (
	"fmt"
	"os/exec"
	"strings"
)

// AptPackageState implements the PackageState interface for the APT package manager
// used by Debian, Ubuntu, and their derivatives.
type AptPackageState struct {
}

func (a *AptPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	// dpkg-query is used instead of "apt list --installed" because it does not
	// require root, has stable machine-readable output, and does not print the
	// "WARNING: apt does not have a stable CLI interface" notice.

	// log the command to be executed
	fmt.Printf("Executing command: dpkg-query -W -f='${Status}' %s\n", name)

	cmd := exec.Command("dpkg-query", "-W", "-f=${Status}", name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Capture stdout and stderr as strings
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		// dpkg-query exits with status 1 if the package is unknown to dpkg
		// https://man7.org/linux/man-pages/man1/dpkg-query.1.html#EXIT_STATUS
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("error checking if package is installed: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	return isInstalledStatus(stdoutStr), nil
}

// isInstalledStatus interprets the dpkg "Status" field for a package.
// A package that was removed but not purged is still known to dpkg (exit
// status 0) with a status such as "deinstall ok config-files", so the
// status text must be checked rather than relying on the exit code alone.
func isInstalledStatus(status string) bool {
	fields := strings.Fields(status)
	if len(fields) < 3 {
		return false
	}
	return fields[len(fields)-1] == "installed"
}

func (a *AptPackageState) IsOutdated(name string, arguments []string) (bool, error) {
	fmt.Printf("Executing command: apt-get -s --only-upgrade install %s\n", name)

	// A simulated install reports which packages would be upgraded without
	// changing the system or requiring root.
	cmd := exec.Command("apt-get", "-s", "--only-upgrade", "install", name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return false, fmt.Errorf("error checking if package is outdated: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	return isOutdatedOutput(stdoutStr, name), nil
}

// isOutdatedOutput reports whether the simulated apt-get output shows an
// upgrade of the given package. apt-get prints a line of the form
// "Inst <name> [<old version>] (<new version> ...)" for each package it
// would upgrade.
func isOutdatedOutput(output string, name string) bool {
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "Inst" && fields[1] == name {
			return true
		}
	}
	return false
}

// InstallCommand returns the command Install runs, without running it.
func (a *AptPackageState) InstallCommand(name string) []string {
	return []string{"sudo", "apt-get", "install", "-y", name}
}

func (a *AptPackageState) Install(name string) error {
	args := a.InstallCommand(name)

	// log the command to be executed
	fmt.Printf("Executing command: %s\n", strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	// Prevent debconf from trying to prompt on a terminal we are not forwarding
	cmd.Env = append(cmd.Environ(), "DEBIAN_FRONTEND=noninteractive")
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

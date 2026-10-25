package snap

import (
	"fmt"
	"os/exec"
	"strings"
)

// SnapPackageState implements the PackageState interface for the snap package
// manager (snapd) used by Ubuntu and available on most other Linux distributions.
type SnapPackageState struct {
}

// notInstalledMessage is what snap prints to stderr (with a non-zero exit
// status) when "snap list <name>" is given a snap that is not installed.
const notInstalledMessage = "no matching snaps installed"

// allUpToDateMessage is what snap prints when "snap refresh --list" finds
// nothing to update.
const allUpToDateMessage = "All snaps up to date."

func (s *SnapPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	// log the command to be executed
	fmt.Printf("Executing command: snap list %s\n", name)

	cmd := exec.Command("snap", "list", name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Capture stdout and stderr as strings
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		// snap exits with a non-zero status both when the snap is not installed
		// and when something else went wrong (snapd not running, etc.), so the
		// message is checked rather than relying on the exit code alone.
		if isNotInstalledMessage(stderrStr) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if package is installed: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	// If we get here, the package is installed
	return true, nil
}

// isNotInstalledMessage reports whether the stderr output of "snap list <name>"
// says the snap is simply not installed, e.g.
// "error: no matching snaps installed".
func isNotInstalledMessage(stderr string) bool {
	return strings.Contains(stderr, notInstalledMessage)
}

func (s *SnapPackageState) IsOutdated(name string, arguments []string) (bool, error) {
	fmt.Printf("Executing command: snap refresh --list\n")

	// "snap refresh --list" shows the snaps that have an update available
	// without changing the system or requiring root.
	cmd := exec.Command("snap", "refresh", "--list")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return false, fmt.Errorf("error checking if package is outdated: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	// Current snapd versions print the "All snaps up to date." message to
	// stderr, leaving stdout empty, but the helper tolerates the message on
	// either stream.
	return isOutdatedOutput(stdoutStr, name), nil
}

// isOutdatedOutput reports whether the output of "snap refresh --list" lists
// the given snap. When updates are available the output is a table with a
// header line followed by one row per snap, where the first column is the
// snap name:
//
//	Name  Version  Rev   Size   Publisher   Notes
//	core  16-2.61  1234  104MB  canonical✓  core
//
// When nothing needs updating snap prints "All snaps up to date." instead.
func isOutdatedOutput(output string, name string) bool {
	for i, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == allUpToDateMessage {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		// skip the table header
		if i == 0 && fields[0] == "Name" {
			continue
		}
		if fields[0] == name {
			return true
		}
	}
	return false
}

// Install installs a snap with "sudo snap install <name>".
//
// Some snaps are published with classic confinement, and snap refuses to
// install them unless --classic is passed explicitly. The PackageState
// interface does not give Install access to the package's extra arguments, so
// there is no way for the metapkg file to request --classic. Mirroring how
// flatpak.go falls back to the flathub remote, the smallest sensible handling
// is to detect snap's "repeat the command including --classic" message and
// retry once with --classic, printing what happened so the user is aware that
// the snap runs outside the usual sandbox.
func (s *SnapPackageState) Install(name string) error {
	args := s.InstallCommand(name)

	// log the command to be executed
	fmt.Printf("Executing command: %s\n", strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Capture stdout and stderr as strings
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		if requiresClassic(stderrStr) {
			classicArgs := classicInstallCommand(name)
			fmt.Printf("Snap %q requires classic confinement. Retrying with --classic...\n", name)
			fmt.Printf("Executing command: %s\n", strings.Join(classicArgs, " "))

			cmd = exec.Command(classicArgs[0], classicArgs[1:]...)
			var classicStdout, classicStderr strings.Builder
			cmd.Stdout = &classicStdout
			cmd.Stderr = &classicStderr
			err = cmd.Run()

			if err == nil {
				return nil // Installation successful with --classic
			}

			// Report the output of the second attempt
			stdoutStr = classicStdout.String()
			stderrStr = classicStderr.String()
		}
		return fmt.Errorf("error installing package: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}
	return nil
}

// InstallCommand returns the command Install runs first, without running it.
func (s *SnapPackageState) InstallCommand(name string) []string {
	return []string{"sudo", "snap", "install", name}
}

// classicInstallCommand returns the command Install retries with when snap
// reports that the snap requires classic confinement.
func classicInstallCommand(name string) []string {
	return []string{"sudo", "snap", "install", "--classic", name}
}

// InstallDescription describes the conditional --classic retry for a dry run.
func (s *SnapPackageState) InstallDescription(name string) []string {
	return []string{
		fmt.Sprintf("may retry with: %s (if snap reports the snap requires classic confinement)",
			strings.Join(classicInstallCommand(name), " ")),
	}
}

// requiresClassic reports whether the stderr output of "snap install <name>"
// says the snap uses classic confinement and must be installed with --classic.
// snap prints a message along the lines of:
//
//	error: This revision of snap "code" was published using classic
//	       confinement and thus may perform arbitrary system changes outside
//	       of the security sandbox that snaps are usually confined to, which
//	       may put your system at risk.
//
//	       If you understand and want to proceed repeat the command including
//	       --classic.
func requiresClassic(stderr string) bool {
	return strings.Contains(stderr, "--classic")
}

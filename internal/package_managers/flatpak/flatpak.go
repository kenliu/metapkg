package flatpak

import (
	"fmt"
	"os/exec"
	"strings"
)

type FlatpakPackageState struct {
}

func (f *FlatpakPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	fmt.Printf("Executing command: flatpak info %s\n", name)

	cmd := exec.Command("flatpak", "info", name)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// flatpak returns exit code 1 if the package is not installed
			if exitError.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, fmt.Errorf("error checking if package is installed: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}

	// If we get here, the package is installed
	return true, nil
}

func (f *FlatpakPackageState) IsOutdated(name string, arguments []string) (bool, error) {
	// Implement this method if needed
	return false, nil
}

// noRemoteChosenMessage is what flatpak prints when more than one remote
// provides the package and none was specified.
const noRemoteChosenMessage = "No remote chosen to resolve matches"

// InstallCommand returns the command Install runs first, without running it.
func (f *FlatpakPackageState) InstallCommand(name string) []string {
	return []string{"flatpak", "install", "-y", name}
}

// flathubInstallCommand returns the command Install retries with when flatpak
// reports that no remote was chosen.
func flathubInstallCommand(name string) []string {
	return []string{"flatpak", "install", "-y", "flathub", name}
}

// InstallDescription describes the conditional flathub retry for a dry run.
func (f *FlatpakPackageState) InstallDescription(name string) []string {
	return []string{
		fmt.Sprintf("may retry with: %s (if flatpak reports no remote chosen)",
			strings.Join(flathubInstallCommand(name), " ")),
	}
}

func (f *FlatpakPackageState) Install(name string) error {
	args := f.InstallCommand(name)
	fmt.Printf("Executing command: %s\n", strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		// Check if the error message indicates no remote was chosen
		if strings.Contains(stderrStr, noRemoteChosenMessage) {
			// If no remote was chosen, try installing with flathub remote
			fmt.Println("No remote specified. Attempting to install from Flathub...")
			fmt.Println("To specify the remote, add the remote name as an additional argument in the metapkg file. For example:")
			fmt.Printf("  %s \"flatpak\" \"fedora\"\n", name)
			flathubArgs := flathubInstallCommand(name)
			fmt.Printf("Executing command: %s\n", strings.Join(flathubArgs, " "))
			cmd = exec.Command(flathubArgs[0], flathubArgs[1:]...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()

			if err == nil {
				return nil // Installation successful with flathub
			}

			// Update stdout and stderr strings after the second attempt
			stdoutStr = stdout.String()
			stderrStr = stderr.String()
		}
		return fmt.Errorf("error installing package: %w\nstdout: %s\nstderr: %s", err, stdoutStr, stderrStr)
	}
	return nil
}

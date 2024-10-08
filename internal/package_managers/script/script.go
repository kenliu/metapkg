package script

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kenliu/metapkg/internal/packages"
)

// ScriptPackageState implements the PackageState interface for script-based installations
type ScriptPackageState struct {
	Scriptdef packages.Scriptdef
}

// NewScriptPackageState creates a new ScriptPackageState with the given scriptdefs
func NewScriptPackageState(scriptdef packages.Scriptdef) *ScriptPackageState {
	state := &ScriptPackageState{
		Scriptdef: scriptdef,
	}
	return state
}

// Checks if a script-based package is installed
func (s *ScriptPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	// For script-based installations, we check if the executable is in the PATH

	// use the `which` command to check if it's installed
	// for now, assume the name of the package is the executable name
	// TODO allow the user to specify the executable name
	return whichCommand(name)
}

// Install runs the script commands to install a package
func (s *ScriptPackageState) Install(name string) error {
	// Create a temporary file for the script
	tmpfile, err := os.CreateTemp("", "install-script-*.sh")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}

	// Ensure the file is both closed and removed at the end of the function
	defer func() {
		tmpfile.Close()
		//TODO create a config flag to not delete the file for debugging purposes
		os.Remove(tmpfile.Name())
	}()

	// Write the script contents to the file
	if len(s.Scriptdef.Commands) == 0 {
		println("warning: no commands found for scriptdef " + s.Scriptdef.Name)
	}

	// print the script contents to the console
	//TODO only print if debug is enabled
	//println("writing script contents:")
	//println(strings.Join(s.Scriptdef.Commands, "\n"))

	scriptContent := "#!/bin/bash\n\n" + strings.Join(s.Scriptdef.Commands, "\n")
	if _, err := tmpfile.Write([]byte(scriptContent)); err != nil {
		return fmt.Errorf("error writing to temp file: %w", err)
	}

	// Close the file immediately after writing
	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("error closing temp file: %w", err)
	}

	// Make the script executable
	if err := os.Chmod(tmpfile.Name(), 0755); err != nil {
		return fmt.Errorf("error making script executable: %w", err)
	}

	// Determine the shell to use
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh" // Fallback to /bin/sh if SHELL is not set
	}

	// Execute the script using the determined shell
	cmd := exec.Command(shell, tmpfile.Name())
	fmt.Printf("Executing command: %v\n", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %v\nOutput: %s", err, string(output))
	}
	// Use the output here
	fmt.Println(string(output))

	if len(output) == 0 {
		fmt.Println("No output captured")
	}

	return nil
}

func (s *ScriptPackageState) IsOutdated(name string, arguments []string) (bool, error) {
	// For script-based installations, we don't have a way to check if a package is outdated
	// so we'll just return false
	return false, nil
}

// whichCommand checks if a command is in the PATH using the `which` command
func whichCommand(name string) (bool, error) {
	// if it returns 0, then it's installed
	// if it returns 1, then it's not installed
	// if it returns an error, then return that error

	println("Checking if " + name + " is installed")
	// 'which' command returns exit status 1 if the executable is not found
	// For any other error, return it
	// If we reach here, the command was successful (exit status 0)
	cmd := exec.Command("which", name)
	err := cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {

			if exitError.ExitCode() == 1 {
				return false, nil
			}
		}

		return false, fmt.Errorf("error checking if package is installed: %w", err)
	}

	return true, nil
}

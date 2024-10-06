package script

import (
	"fmt"
	"os/exec"

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

	// assume the name of the package is the executable name
	// TODO allow the user to specify the executable name
	// use the `which` command to check if it's installed

	// if it returns 0, then it's installed
	// if it returns 1, then it's not installed
	// if it returns an error, then return that error

	println("Checking if " + name + " is installed")
	cmd := exec.Command("which", name)
	err := cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// 'which' command returns exit status 1 if the executable is not found
			if exitError.ExitCode() == 1 {
				return false, nil
			}
		}
		// For any other error, return it
		return false, fmt.Errorf("error checking if package is installed: %w", err)
	}

	// If we reach here, the command was successful (exit status 0)
	return true, nil
}

// Install runs the script commands to install a package
func (s *ScriptPackageState) Install(name string) error {
	println("Script package manager not implemented yet")
	// for _, cmd := range s.Scriptdef.Commands {
	// 	fmt.Printf("Executing command: %s\n", cmd)

	// 	// Split the command into parts
	// 	parts := strings.Fields(cmd)
	// 	if len(parts) == 0 {
	// 		return fmt.Errorf("empty command in scriptdef: %s", name)
	// 	}

	// 	// Create the command
	// 	execCmd := exec.Command(parts[0], parts[1:]...)

	// 	// Run the command
	// 	output, err := execCmd.CombinedOutput()
	// 	if err != nil {
	// 		return fmt.Errorf("error executing command '%s': %w\nOutput: %s", cmd, err, string(output))
	// 	}

	// 	fmt.Printf("Command output:\n%s\n", string(output))
	//}

	return nil
}

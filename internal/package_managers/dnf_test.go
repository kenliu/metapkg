package package_managers

import (
	"os"
	"os/exec"
	"testing"
)

// mockCmd is a helper function to mock exec.Command
func mockCmd(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

//func TestDnfPackageState_IsInstalled(t *testing.T) {
//	tests := []struct {
//		name           string
//		packageName    string
//		mockStdout     string
//		mockStderr     string
//		mockExitStatus int
//		expected       bool
//		expectError    bool
//	}{
//		{
//			name:           "Package is installed",
//			packageName:    "test-package",
//			mockStdout:     "Installed Packages\ntest-package.x86_64 1.2.3-4 @System",
//			mockStderr:     "",
//			mockExitStatus: 0,
//			expected:       true,
//			expectError:    false,
//		},
//		{
//			name:           "Package is not installed",
//			packageName:    "nonexistent-package",
//			mockStdout:     "",
//			mockStderr:     "Error: No matching Packages to list",
//			mockExitStatus: 1,
//			expected:       false,
//			expectError:    false,
//		},
//		{
//			name:           "DNF command fails",
//			packageName:    "error-package",
//			mockStdout:     "",
//			mockStderr:     "Error: DNF command failed",
//			mockExitStatus: 1,
//			expected:       false,
//			expectError:    true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// Mock exec.Command
//			execCommand = mockCmd
//			defer func() { execCommand = exec.Command }()
//
//			// Set up the mock command behavior
//			mockExecCommand = func(command string, args ...string) *exec.Cmd {
//				cmd := mockCmd(command, args...)
//				cmd.Stdout.Write([]byte(tt.mockStdout))
//				cmd.Stderr.Write([]byte(tt.mockStderr))
//				return cmd
//			}
//
//			d := &DnfPackageState{}
//			result, err := d.IsInstalled(tt.packageName, nil)
//
//			if tt.expectError {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//			assert.Equal(t, tt.expected, result)
//		})
//	}
//}

// TestHelperProcess isn't a real test. It's used to mock exec.Command in the tests above.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

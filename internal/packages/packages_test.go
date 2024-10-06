package packages

import (
	"strings"
	"testing"

	"github.com/sblinch/kdl-go"
)

func TestParseDnf(t *testing.T) {
	// Define a small KDL snippet
	kdlString := `
foo "dnf"
`

	// Parse the KDL string
	doc, err := kdl.Parse(strings.NewReader(kdlString))
	if err != nil {
		t.Fatalf("Failed to parse KDL: %v", err)
	}

	// Get the first (and only) node
	node := doc.Nodes[0]

	// Call the function being tested
	pkg, err := parseDnf(node)

	// Assert the results
	if err != nil {
		t.Fatalf("parseDnf returned an error: %v", err)
	}

	// Check the parsed package details
	if pkg.Name != "foo" {
		t.Errorf("Expected package name 'foo', got '%s'", pkg.Name)
	}

	if pkg.PackageManager != "dnf" {
		t.Errorf("Expected package manager 'dnf', got '%s'", pkg.PackageManager)
	}
}

func TestParseFlatpak(t *testing.T) {
	// Define a small KDL snippet
	kdlString := `
foo "flatpak"
`
	// Parse the KDL string
	doc, err := kdl.Parse(strings.NewReader(kdlString))
	if err != nil {
		t.Fatalf("Failed to parse KDL: %v", err)
	}

	// Get the first (and only) node
	node := doc.Nodes[0]

	// Call the function being tested
	pkg, err := parseFlatpak(node)

	// Assert the results
	if err != nil {
		t.Fatalf("parseFlatpak returned an error: %v", err)
	}

	// Check the parsed package details
	if pkg.Name != "foo" {
		t.Errorf("Expected package name 'foo', got '%s'", pkg.Name)
	}

	if pkg.PackageManager != "flatpak" {
		t.Errorf("Expected package manager 'flatpak', got '%s'", pkg.PackageManager)
	}
}

func TestParseScript(t *testing.T) {
	// Define a small KDL snippet
	kdlString := `
foo "script"
`

	// Parse the KDL string
	doc, err := kdl.Parse(strings.NewReader(kdlString))
	if err != nil {
		t.Fatalf("Failed to parse KDL: %v", err)
	}

	// Get the first (and only) node
	node := doc.Nodes[0]

	// Call the function being tested
	pkg, err := parseScript(node)

	// Assert the results
	if err != nil {
		t.Fatalf("parseScript returned an error: %v", err)
	}

	// Check the parsed package details
	if pkg.Name != "foo" {
		t.Errorf("Expected package name 'foo', got '%s'", pkg.Name)
	}

	if pkg.PackageManager != "script" {
		t.Errorf("Expected package manager 'script', got '%s'", pkg.PackageManager)
	}
}

func TestParseScriptdef(t *testing.T) {
	// Define a KDL snippet with various command types
	kdlString := `
scriptdef "test_script" {
	cmd "echo 'hello'"
	"ls -la"
	cmd "pwd"
}
`

	// Parse the KDL string
	doc, err := kdl.Parse(strings.NewReader(kdlString))
	if err != nil {
		t.Fatalf("Failed to parse KDL: %v", err)
	}

	// Get the first (and only) node
	node := doc.Nodes[0]

	// Call the function being tested
	scriptdef, err := parseScriptdef(node)

	// Assert the results
	if err != nil {
		t.Fatalf("parseScriptdef returned an error: %v", err)
	}

	// Check the parsed scriptdef details
	if scriptdef.Name != "test_script" {
		t.Errorf("Expected scriptdef name 'test_script', got '%s'", scriptdef.Name)
	}

	expectedCommands := []string{
		"echo 'hello'",
		"ls -la",
		"pwd",
	}

	if len(scriptdef.Commands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(scriptdef.Commands))
	}

	for i, cmd := range scriptdef.Commands {
		if cmd != expectedCommands[i] {
			t.Errorf("Expected command '%s', got '%s'", expectedCommands[i], cmd)
		}
	}
}

func TestParseDefaultPackageManager(t *testing.T) {
	// Define a small KDL snippet
	kdlString := `
foo
`

	// Parse the KDL string
	doc, err := kdl.Parse(strings.NewReader(kdlString))
	if err != nil {
					t.Fatalf("Failed to parse KDL: %v", err)
	}

	// Get the first (and only) node
	node := doc.Nodes[0]

	// Call the function being tested
	pkg, err := parseDefaultPackageManager(node)

	// Assert the results
	if err != nil {
		t.Fatalf("parseDefaultPackageManager returned an error: %v", err)
	}

	// Check the parsed package details
	if pkg.Name != "foo" {
		t.Errorf("Expected package name 'example', got '%s'", pkg.Name)
	}

	if pkg.PackageManager != "dnf" {
		t.Errorf("Expected package manager 'dnf', got '%s'", pkg.PackageManager)
	}
}

func TestLoadMetapackage(t *testing.T) {
	metapackage, err := LoadMetapackageFile("../../test-files/sample_config.kdl")
	if err != nil {
		t.Fatalf("Failed to load metapackage: %v", err)
	}

	// Check the parsed metapackage details
	if len(metapackage.Packages) != 16 {
		t.Errorf("Expected 4 packages, got %d", len(metapackage.Packages))
	}

	if len(metapackage.Scriptdefs) != 2 {
		t.Errorf("Expected 2 scripts, got %d", len(metapackage.Scriptdefs))
	}
}

package packages

import (
	"github.com/sblinch/kdl-go"
	"strings"
	"testing"
)

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

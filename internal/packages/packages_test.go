package packages

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/njreid/gokdl2"
)

func TestParseDnf(t *testing.T) {
	// Define a small KDL snippet
	kdlString := `
foo "dnf"
`

	// Parse the KDL string
	doc, err := kdl.ParseWithOptions(strings.NewReader(kdlString), kdl.ParseOptions{Version: kdl.ParseVersionV2})
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

func TestParseApt(t *testing.T) {
	kdlString := `
foo "apt"
`
	doc, err := kdl.ParseWithOptions(strings.NewReader(kdlString), kdl.ParseOptions{Version: kdl.ParseVersionV2})
	if err != nil {
		t.Fatalf("Failed to parse KDL: %v", err)
	}

	pkg, err := parseApt(doc.Nodes[0])
	if err != nil {
		t.Fatalf("parseApt returned an error: %v", err)
	}

	if pkg.Name != "foo" {
		t.Errorf("Expected package name 'foo', got '%s'", pkg.Name)
	}

	if pkg.PackageManager != "apt" {
		t.Errorf("Expected package manager 'apt', got '%s'", pkg.PackageManager)
	}
}

func TestParseFlatpak(t *testing.T) {
	// Define a small KDL snippet
	kdlString := `
foo "flatpak"
`
	// Parse the KDL string
	doc, err := kdl.ParseWithOptions(strings.NewReader(kdlString), kdl.ParseOptions{Version: kdl.ParseVersionV2})
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

func TestParseSnap(t *testing.T) {
	kdlString := `
foo "snap"
`
	doc, err := kdl.ParseWithOptions(strings.NewReader(kdlString), kdl.ParseOptions{Version: kdl.ParseVersionV2})
	if err != nil {
		t.Fatalf("Failed to parse KDL: %v", err)
	}

	pkg, err := parseSnap(doc.Nodes[0])
	if err != nil {
		t.Fatalf("parseSnap returned an error: %v", err)
	}

	if pkg.Name != "foo" {
		t.Errorf("Expected package name 'foo', got '%s'", pkg.Name)
	}

	if pkg.PackageManager != "snap" {
		t.Errorf("Expected package manager 'snap', got '%s'", pkg.PackageManager)
	}
}

func TestParseScript(t *testing.T) {
	// Define a small KDL snippet
	kdlString := `
foo "script"
`

	// Parse the KDL string
	doc, err := kdl.ParseWithOptions(strings.NewReader(kdlString), kdl.ParseOptions{Version: kdl.ParseVersionV2})
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
	doc, err := kdl.ParseWithOptions(strings.NewReader(kdlString), kdl.ParseOptions{Version: kdl.ParseVersionV2})
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
	doc, err := kdl.ParseWithOptions(strings.NewReader(kdlString), kdl.ParseOptions{Version: kdl.ParseVersionV2})
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

	// The default depends on the operating system the test runs on
	expected := DefaultPackageManager()
	if pkg.PackageManager != expected {
		t.Errorf("Expected package manager '%s', got '%s'", expected, pkg.PackageManager)
	}
}

func TestLoadMetapackage(t *testing.T) {
	metapackage, err := LoadMetapackageFile("../../test-files/sample_config.kdl")
	if err != nil {
		t.Fatalf("Failed to load metapackage: %v", err)
	}

	// Check the parsed metapackage details
	if len(metapackage.Packages) != 8 {
		t.Errorf("Expected 8 packages, got %d", len(metapackage.Packages))
	}

	if len(metapackage.Scriptdefs) != 2 {
		t.Errorf("Expected 2 scripts, got %d", len(metapackage.Scriptdefs))
	}
	commands := metapackage.Scriptdefs["gcloud"].Commands
	if len(commands) != 2 {
		t.Fatalf("Expected 2 gcloud commands, got %d", len(commands))
	}
	if !strings.HasPrefix(commands[0], "sudo tee -a ") || !strings.HasSuffix(commands[0], "\nEOM") {
		t.Errorf("Multiline command was not preserved: %q", commands[0])
	}
}

func TestLoadMetapackageKDLV2(t *testing.T) {
	file := filepath.Join(t.TempDir(), "metapkg.kdl")
	input := `
"default package"
"dnf package" dnf
org.kde.filelight flatpak
"brew package" brew
custom script
scriptdef custom {
    cmd #"echo "hello""#
    #"""
        printf '%s\n' "world"
        echo done
        """#
}
`
	if err := os.WriteFile(file, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadMetapackageFile(file)
	if err != nil {
		t.Fatal(err)
	}
	want := &Metapackage{
		Packages: []Package{
			{Name: "default package", PackageManager: DefaultPackageManager()},
			{Name: "dnf package", PackageManager: "dnf"},
			{Name: "org.kde.filelight", PackageManager: "flatpak"},
			{Name: "brew package", PackageManager: "brew"},
			{Name: "custom", PackageManager: "script"},
		},
		Scriptdefs: map[string]Scriptdef{
			"custom": {
				Name:     "custom",
				Commands: []string{`echo "hello"`, "printf '%s\\n' \"world\"\necho done"},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Loaded configuration = %#v, want %#v; scripts = %#v", got, want, got.Scriptdefs)
	}
}

func TestLoadMetapackageRejectsKDLV1(t *testing.T) {
	file := filepath.Join(t.TempDir(), "metapkg.kdl")
	if err := os.WriteFile(file, []byte("foo r\"dnf\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadMetapackageFile(file); err == nil || !strings.Contains(err.Error(), "error parsing KDL v2") {
		t.Fatalf("Expected a KDL v2 parse error, got %v", err)
	}
}

func TestExampleFilesUseKDLV2(t *testing.T) {
	for _, pattern := range []string{"../../test-files/*.kdl", "../../experiments/*.kdl"} {
		files, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			t.Run(file, func(t *testing.T) {
				data, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := kdl.ParseWithOptions(strings.NewReader(string(data)), kdl.ParseOptions{Version: kdl.ParseVersionV2}); err != nil {
					t.Fatalf("Example is not valid KDL v2: %v", err)
				}
			})
		}
	}
}

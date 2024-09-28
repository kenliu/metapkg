package packages

import (
	"fmt"
	"os"

	"github.com/sblinch/kdl-go"
	"github.com/sblinch/kdl-go/document"
)

// TODO: make this configurable
const DEFAULT_PACKAGE_MANAGER = "dnf"

type Package struct {
	Name           string
	PackageManager string
	Arguments      []string
}

type Scriptdef struct {
	Name     string
	Commands []string
}

type Metapackage struct {
	Packages   []Package
	Scriptdefs map[string]Scriptdef
}

func LoadMetapackage(file string) (*Metapackage, error) {
	// Open the file
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Parse the KDL document
	doc, err := kdl.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("error parsing KDL: %w", err)
	}

	metapackage := &Metapackage{
		Scriptdefs: make(map[string]Scriptdef),
	}

	// Iterate through the root nodes
	for _, node := range doc.Nodes {
		switch node.Name.String() {
		case "flatpak":
			println("found flatpak node")
			break
		case "dnf":
			println("found dnf node")
			pkg, err := parseDnf(node)
			if err != nil {
				return nil, fmt.Errorf("error parsing dnf: %w", err)
			}
			metapackage.Packages = append(metapackage.Packages, *pkg)
		case "script":
			println("found script node")
			break
		case "scriptdef":
			println("found scriptdef node")
			scriptdef, err := parseScriptdef(node)
			if err != nil {
				return nil, fmt.Errorf("error parsing scriptdef: %w", err)
			}
			metapackage.Scriptdefs[scriptdef.Name] = *scriptdef
		default:
			pkg, err := parseDefaultPackageManager(node)
			if err != nil {
				return nil, fmt.Errorf("error parsing default package manager: %w", err)
			}
			metapackage.Packages = append(metapackage.Packages, *pkg)
		}
	}

	return metapackage, nil
}

func parseDnf(node *document.Node) (*Package, error) {
	pkg := &Package{}
	pkg.PackageManager = "dnf"
	pkg.Name = node.Arguments[0].ValueString()
	return pkg, nil
}

func parseFlatpak(node *document.Node) (*Package, error) {
	pkg := &Package{}
	pkg.PackageManager = "flatpak"
	pkg.Name = node.Arguments[0].ValueString()
	return pkg, nil
}

func parseScript(node *document.Node) (*Package, error) {
	pkg := &Package{}
	pkg.PackageManager = "script"
	pkg.Name = node.Arguments[0].ValueString()
	return pkg, nil
}

func parseDefaultPackageManager(node *document.Node) (*Package, error) {
	//parse the case where there is a node without a named package manager
	//the default package manager with no options is used
	pkg := &Package{}

	// The package manager is not specified, so the default is used
	pkg.PackageManager = DEFAULT_PACKAGE_MANAGER
	pkg.Name = node.String()

	return pkg, nil
}

func parseScriptdef(node *document.Node) (*Scriptdef, error) {
	scriptdef := &Scriptdef{
		Name: node.Arguments[0].ValueString(),
	}
	for _, child := range node.Children {
		if child.Name.String() == "cmd" {
			if len(child.Arguments) > 0 {
				scriptdef.Commands = append(scriptdef.Commands, child.Arguments[0].ValueString())
			}
		}
	}
	return scriptdef, nil
}

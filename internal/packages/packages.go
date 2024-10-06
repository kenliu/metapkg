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

func LoadMetapackageFile(file string) (*Metapackage, error) {
	//uncomment this to test with the example metapkg file while debugging in VSCode
	//file = "/home/kenliu/code/metapkg/bin/metapkg.kdl"
	println("Loading metapkg file")
	println("====================")
	println("Loading file: " + file)

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

	// TODO handle some basic validation of the document
	// for example that there is no duplicate scriptdef names
	// and that there is no duplicate package names
	// and that each script has a corresponding scriptdef

	metapackage := &Metapackage{
		Scriptdefs: make(map[string]Scriptdef),
	}

	// Iterate through the root nodes
	for _, node := range doc.Nodes {
		switch node.Name.String() {
		case "scriptdef":
			println("found scriptdef node")
			scriptdef, err := parseScriptdef(node)
			if err != nil {
				return nil, fmt.Errorf("error parsing scriptdef: %w", err)
			}
			metapackage.Scriptdefs[scriptdef.Name] = *scriptdef
		default:
			pkg := &Package{}

			// we know the node is a package manager because it's not a script or scriptdef
			// so we can parse the first argument as the package manager name
			// and the nodename as the package name
			// and the rest as the arguments
			// then pass the node to the appropriate parse function that
			// corresponds to the package manager

			// if there are no arguments, then it's the case where there is no package manager specified
			// we pass the node to the default package manager parser
			// which will parse the node as a package manager with no options

			if len(node.Arguments) == 0 {
				pkg, err := parseDefaultPackageManager(node)
				if err != nil {
					return nil, fmt.Errorf("error parsing default package manager: %w", err)
				}
				metapackage.Packages = append(metapackage.Packages, *pkg)
			} else {
				packageManagerName := node.Arguments[0].ValueString()
				switch packageManagerName {
				case "script":
					println("found script node")
					pkg, err := parseScript(node)
					if err != nil {
						return nil, fmt.Errorf("error parsing script: %w", err)
					}
					metapackage.Packages = append(metapackage.Packages, *pkg)
				case "dnf":
					println("found dnf node")
					pkg, err := parseDnf(node)
					if err != nil {
						return nil, fmt.Errorf("error parsing dnf: %w", err)
					}
					metapackage.Packages = append(metapackage.Packages, *pkg)
				case "flatpak":
					println("found flatpak node")
					pkg, err := parseFlatpak(node)
					if err != nil {
						return nil, fmt.Errorf("error parsing flatpak: %w", err)
					}
					metapackage.Packages = append(metapackage.Packages, *pkg)
				case "brew":
					println("found brew node")
					pkg, err := parseBrew(node)
					if err != nil {
						return nil, fmt.Errorf("error parsing brew: %w", err)
					}
					metapackage.Packages = append(metapackage.Packages, *pkg)
				default:
					println("unknown package manager: %s", packageManagerName)
					return nil, fmt.Errorf("unknown package manager: %s", packageManagerName)
				}
			}

			metapackage.Packages = append(metapackage.Packages, *pkg)
		}
	}

	return metapackage, nil
}

func parseDnf(node *document.Node) (*Package, error) {
	pkg := &Package{}
	pkg.PackageManager = "dnf"
	pkg.Name = node.Name.String()
	return pkg, nil
}

func parseFlatpak(node *document.Node) (*Package, error) {
	pkg := &Package{}
	pkg.PackageManager = "flatpak"
	pkg.Name = node.Name.String()
	return pkg, nil
}

// parse a "brew" package manager node
func parseBrew(node *document.Node) (*Package, error) {
	pkg := &Package{}
	pkg.PackageManager = "brew"
	pkg.Name = node.Name.String()
	return pkg, nil
}

func parseScript(node *document.Node) (*Package, error) {
	pkg := &Package{}
	pkg.PackageManager = "script"
	pkg.Name = node.Name.String()
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
		} else {
			scriptdef.Commands = append(scriptdef.Commands, child.Name.ValueString())
		}
	}
	return scriptdef, nil
}

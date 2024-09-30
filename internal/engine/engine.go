package engine

import (
	"fmt"
	"github.com/kenliu/metapkg/internal/package_managers/dnf"

	"github.com/kenliu/metapkg/internal/packages"
)

func InstallPackages(m *packages.Metapackage) error {
	//iterate over each package, checking if each one is installed first
	for _, pkg := range m.Packages {
		// if a package is not installed, install it
		if pkg.PackageManager == "dnf" {
			dnf := dnf.DnfPackageState{}
			installed, err := dnf.IsInstalled(pkg.Name, pkg.Arguments)
			if err != nil {
				return err
			}
			if !installed {
				dnf.Install(pkg.Name)
			}
		}
	}
	return nil
}

func ListOutdatedPackages(m *packages.Metapackage) error {
	// Implementation for listing outdated packages
	// ...
	return fmt.Errorf("not implemented")
}

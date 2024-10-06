package engine

import (
	"fmt"
	"github.com/kenliu/metapkg/internal/package_managers/dnf"
	"github.com/kenliu/metapkg/internal/package_managers/flatpak"
	"github.com/kenliu/metapkg/internal/package_managers/script"

	"github.com/kenliu/metapkg/internal/packages"
)

func InstallPackages(m *packages.Metapackage) error {
	println()
	println("Installing packages")
	println("===================")
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
					err = dnf.Install(pkg.Name)
				if err != nil {
					return err
				}
			}
		} else if pkg.PackageManager == "flatpak" {
			flatpak := flatpak.FlatpakPackageState{}
			installed, err := flatpak.IsInstalled(pkg.Name, pkg.Arguments)
			if err != nil {
				return err
			}
			if !installed {
				err = flatpak.Install(pkg.Name)
				if err != nil {
					return err
				}
			}
		} else if pkg.PackageManager == "script" {
			script := script.ScriptPackageState{}
			installed, err := script.IsInstalled(pkg.Name, pkg.Arguments)
			if err != nil {
				return err
			}
			if !installed {
				err = script.Install(pkg.Name)
				if err != nil {
					return err
				}
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

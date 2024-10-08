package engine

import (
	"fmt"
	"github.com/kenliu/metapkg/internal/package_managers/dnf"
	"github.com/kenliu/metapkg/internal/package_managers/flatpak"
	"github.com/kenliu/metapkg/internal/package_managers/script"

	"github.com/kenliu/metapkg/internal/packages"
	"github.com/kenliu/metapkg/internal/package_managers"
)

func InstallPackages(m *packages.Metapackage) error {
	println("\nInstalling packages")
	println("===================")

	for _, pkg := range m.Packages {
		var packageState package_managers.PackageState

		switch pkg.PackageManager {
		case "dnf":
			packageState = &dnf.DnfPackageState{}
		case "flatpak":
			packageState = &flatpak.FlatpakPackageState{}
		case "script":
			packageState = script.NewScriptPackageState(m.Scriptdefs[pkg.Name])
		default:
			return fmt.Errorf("unsupported package manager: %s", pkg.PackageManager)
		}

		installed, err := packageState.IsInstalled(pkg.Name, pkg.Arguments)
		if err != nil {
			return err
		}

		if !installed {
			if err := packageState.Install(pkg.Name); err != nil {
				return err
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

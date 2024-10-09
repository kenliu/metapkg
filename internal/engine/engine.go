package engine

import (
	"fmt"
	"github.com/kenliu/metapkg/internal/package_managers/brew"
	"github.com/kenliu/metapkg/internal/package_managers/dnf"
	"github.com/kenliu/metapkg/internal/package_managers/flatpak"
	"github.com/kenliu/metapkg/internal/package_managers/script"

	"github.com/kenliu/metapkg/internal/packages"
	"github.com/kenliu/metapkg/internal/package_managers"
)

// packageManagerFactory is a map of package manager names to their respective PackageState implementations
var packageManagerFactory = map[string]func() package_managers.PackageState{
	"dnf":     func() package_managers.PackageState { return &dnf.DnfPackageState{} },
	"flatpak": func() package_managers.PackageState { return &flatpak.FlatpakPackageState{} },
	"brew":    func() package_managers.PackageState { return &brew.BrewPackageState{} },
}

func InstallPackages(m *packages.Metapackage) error {
	println("\nInstalling packages")
	println("===================")

	for _, pkg := range m.Packages {
		var packageState package_managers.PackageState

		if pkg.PackageManager == "script" {
			packageState = script.NewScriptPackageState(m.Scriptdefs[pkg.Name])
		} else {
			stateFactory, ok := packageManagerFactory[pkg.PackageManager]
			if !ok {
				return fmt.Errorf("unsupported package manager: %s", pkg.PackageManager)
			}
			packageState = stateFactory()
		}

		if err := handlePackage(packageState, pkg); err != nil {
			return err
		}
	}

	return nil
}

func handlePackage(packageState package_managers.PackageState, pkg packages.Package) error {
	installed, err := packageState.IsInstalled(pkg.Name, pkg.Arguments)
	if err != nil {
		return err
	}

	if !installed {
		if err := packageState.Install(pkg.Name); err != nil {
			return err
		}
	}

	return nil
}

func ListOutdatedPackages(m *packages.Metapackage) error {
	// Implementation for listing outdated packages
	// ...
	return fmt.Errorf("not implemented")
}

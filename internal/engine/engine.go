package engine

import (
	"fmt"
	"io"
	"strings"

	"github.com/kenliu/metapkg/internal/package_managers/apt"
	"github.com/kenliu/metapkg/internal/package_managers/brew"
	"github.com/kenliu/metapkg/internal/package_managers/dnf"
	"github.com/kenliu/metapkg/internal/package_managers/flatpak"
	"github.com/kenliu/metapkg/internal/package_managers/script"
	"github.com/kenliu/metapkg/internal/package_managers/snap"

	"github.com/kenliu/metapkg/internal/package_managers"
	"github.com/kenliu/metapkg/internal/packages"
)

// packageManagerFactory is a map of package manager names to their respective PackageState implementations
var packageManagerFactory = map[string]func() package_managers.PackageState{
	"dnf":     func() package_managers.PackageState { return &dnf.DnfPackageState{} },
	"apt":     func() package_managers.PackageState { return &apt.AptPackageState{} },
	"flatpak": func() package_managers.PackageState { return &flatpak.FlatpakPackageState{} },
	"brew":    func() package_managers.PackageState { return &brew.BrewPackageState{} },
	"snap":    func() package_managers.PackageState { return &snap.SnapPackageState{} },
}

// newPackageState returns the PackageState that manages the given package.
// It is a variable so tests can substitute fakes.
var newPackageState = func(m *packages.Metapackage, pkg packages.Package) (package_managers.PackageState, error) {
	if pkg.PackageManager == "script" {
		return script.NewScriptPackageState(m.Scriptdefs[pkg.Name]), nil
	}

	stateFactory, ok := packageManagerFactory[pkg.PackageManager]
	if !ok {
		return nil, fmt.Errorf("unsupported package manager: %s", pkg.PackageManager)
	}
	return stateFactory(), nil
}

func InstallPackages(m *packages.Metapackage) error {
	println("\nInstalling packages")
	println("===================")

	for _, pkg := range m.Packages {
		packageState, err := newPackageState(m, pkg)
		if err != nil {
			return err
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

// planEntry is one line of the dry-run plan: a package, whether it is already
// installed, and (if not) what Install would do.
type planEntry struct {
	pkg       packages.Package
	installed bool
	// command is the single command Install would run, or nil if the package
	// manager does not describe one
	command []string
	// details are extra lines describing conditional behaviour or script
	// commands, or nil
	details []string
}

// DryRunInstall reports what InstallPackages would do for the metapackage
// without changing the system. Only read-only checks (IsInstalled) are run;
// Install is never called. The plan is written to out, one line per package
// in file order, followed by a summary.
func DryRunInstall(m *packages.Metapackage, out io.Writer) error {
	// Run all the checks first so the plan prints as one block, separate from
	// the "Executing command: ..." lines the checks log.
	var entries []planEntry
	for _, pkg := range m.Packages {
		packageState, err := newPackageState(m, pkg)
		if err != nil {
			return err
		}

		entry, err := planPackage(packageState, pkg)
		if err != nil {
			return err
		}
		entries = append(entries, entry)
	}

	printPlan(out, entries)
	return nil
}

func planPackage(packageState package_managers.PackageState, pkg packages.Package) (planEntry, error) {
	entry := planEntry{pkg: pkg}

	installed, err := packageState.IsInstalled(pkg.Name, pkg.Arguments)
	if err != nil {
		return entry, err
	}
	entry.installed = installed
	if installed {
		return entry, nil
	}

	if commander, ok := packageState.(package_managers.InstallCommander); ok {
		entry.command = commander.InstallCommand(pkg.Name)
	}
	if describer, ok := packageState.(package_managers.InstallDescriber); ok {
		entry.details = describer.InstallDescription(pkg.Name)
	}
	return entry, nil
}

func printPlan(out io.Writer, entries []planEntry) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Dry run: no changes will be made")
	fmt.Fprintln(out, "================================")

	toInstall, alreadyInstalled := 0, 0
	for _, entry := range entries {
		if entry.installed {
			alreadyInstalled++
			fmt.Fprintf(out, "[skip] %s (%s): already installed\n", entry.pkg.Name, entry.pkg.PackageManager)
			continue
		}
		toInstall++

		switch {
		case entry.command != nil:
			fmt.Fprintf(out, "[install] %s (%s): would run: %s\n", entry.pkg.Name, entry.pkg.PackageManager, strings.Join(entry.command, " "))
		case entry.details != nil:
			fmt.Fprintf(out, "[install] %s (%s): would run:\n", entry.pkg.Name, entry.pkg.PackageManager)
		default:
			fmt.Fprintf(out, "[install] %s (%s): would install\n", entry.pkg.Name, entry.pkg.PackageManager)
		}
		for _, line := range entry.details {
			fmt.Fprintf(out, "    %s\n", line)
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, "Summary: %d to install, %d already installed\n", toInstall, alreadyInstalled)
	fmt.Fprintln(out, "Dry run: nothing was changed")
}

func ListOutdatedPackages(m *packages.Metapackage) error {
	// Implementation for listing outdated packages
	// ...
	return fmt.Errorf("not implemented")
}

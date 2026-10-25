package package_managers

// PackageState defines the interface for managing a package across different package management systems
type PackageState interface {
	// IsInstalled checks if a package is installed
	// package: the name of the package to check
	// arguments: additional arguments for the package manager
	// returns: true if the package is installed, false otherwise, and any error encountered
	IsInstalled(name string, arguments []string) (bool, error)

	// IsOutdated checks if a package is outdated
	// package: the name of the package to check
	// arguments: additional arguments for the package manager
	// returns: true if the package is outdated, false otherwise, and any error encountered
	IsOutdated(name string, arguments []string) (bool, error)

	// Install installs a package
	// package: the name of the package to install
	// returns: any error encountered during installation
	Install(name string) error
}

// InstallCommander is an optional interface for a PackageState that installs
// a package by running a single command. InstallCommand returns that command
// (the program followed by its arguments) without running it, so a dry run
// can report what Install would do. Install must build its command from
// InstallCommand so the reported and executed commands cannot drift apart.
type InstallCommander interface {
	// InstallCommand returns the command Install runs for the given package
	InstallCommand(name string) []string
}

// InstallDescriber is an optional interface for a PackageState whose Install
// does more than run the single command reported by InstallCommand: for
// example a conditional retry with different flags, or a generated shell
// script. InstallDescription returns extra lines for the dry-run plan that
// describe that behaviour, one entry per line.
type InstallDescriber interface {
	// InstallDescription returns lines describing what Install does for the
	// given package beyond (or instead of) InstallCommand
	InstallDescription(name string) []string
}

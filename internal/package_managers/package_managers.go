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


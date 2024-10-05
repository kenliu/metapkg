# Overview
A command-line go program called "metapackage" using the specification in this file that will keep track of a user's installed dnf packages in a file called metapackage.kdl and run package manager commands and shell scriptlets to install, remove, and update packages based on the configuration given in the metapackage.kdl file.

# metapackage-conf.kdl configuration file
- The "metapackage.conf" file is a KDL-format file that contains the user's configuration for running the metapackage program. 
- The program should look for the file in the current directory first, then in the user's home directory.
- If this file is not found, the program should use a set of default values for the configuration.
- The metapackage.conf file should be able to specify which package manager to use by default. The key is called "default_package_manager" and the value is the name of the package manager to use. The possible values for this option are "dnf", "flatpak", "snap", and "brew".

# metapackage.kdl file specification
- The program should look in the current directory for a file named "metapackage.kdl" and use that to determine which packages should be installed.
- The first set of keys in the metapackage.kdl file are the packages to install.
- The second set of keys in the metapackage.kdl file are the scriptlets to run. These are all nodes with a "script" key, followed by a string that is the name of the scriptlet.
- comments are permitted in the metapackage file, as defined by the KDL file specification


## Package nodes
- The package nodes can have any arbitrary string as the name of the node.
- There are a set of reserved keys that may not be used by package nodes. These are: (currently none)
- If a package node has a key that matches the name of one of the supported package managers, then the package node should have a second argument with a string that is the name of the package to install.
- If a package node has a key that is not the name of a supported package manager, then the default package manager is used and the node name is the package name.
- If a package node has a key called "script" then the package node should have a second argument with a string that is the name of the scriptlet to run, matching a key of one of the scriptlet nodes. The second argument must match the name of a scriptlet node.
- some package manager names have short aliases defined for convenience:
    - flatpak: fp
    - script: s
    - dnf: yum

Here is a sample of some package nodes, showing the case of package names as the first node and another case of the node being a package manager name followed by the name of a package:

```
bat
eza
flatpak "org.kde.filelight"
script "vivaldi"
```


## Scriptdef nodes
- These are all nodes with a "scriptdef" key, followed by a string that is the name of the scriptdef.
- The scriptdef nodes can have any number of child nodes. These represent the shell commands that make up the scriptdef.
- The scriptdef nodes can simply be strings that are executed as shell commands.
- By default the scriptdef nodes are executed in a subshell with the shell set to value of SHELL environment variable.

Here is a sample of a scriptdef node showing the "scriptdef" node name followed by the name of the scriptdef:

```
scriptdef "vivaldi" {
  "sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo"
}
```

# metapackage command line interface
- The program should have a "metapackage" command that the user can run.
- "install" subcommand that the user can run to install the packages specified in the metapackage.kdl file.
- "outdated" subcommand that the user can run to list all packages that are present in the metapackage.kdl file that are out of date and need to be updated.
- A "version" subcommand that prints out the version of the program and exits.
- A "help" subcommand that prints out the help message and exits.
- A "dry run" global flag that the user can run to print out the commands that would be run without actually running them.
- A "file" global flag that allows the user to specify a different file to use for the configuration.
- A "debug" global flag that prints out the debug information while running.
- A "verbose" global flag that prints out the verbose information while running.
- A "quiet" global flag that suppresses all output while running.

# metapackage execution specification
- The program should use the dnf package manager to install and update packages by default.
- The program also supports installing packages using the apt, flatpak, snap, and brew package managers.
- The program should also allow the user to specify in the metapackage.conf.kdl file which package manager to use by default.
- The program should also support running arbitrary shell commands using shell scripts to install packages.
- The program should also support a "dry run" mode where it prints out the commands it would run without actually running them.

# other requirements
- Also generate a sample metapackage.kdl file that could hypothetically be used as input for this program.
- Generate unit tests for the program that run in the golang test framework.
- Generate a Dockerfile for running the program in a container to run automated integration tests.
- Generate a README.md file with usage instructions for the program.
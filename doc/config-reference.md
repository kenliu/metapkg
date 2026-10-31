# Config reference

A metapkg file is written in [KDL v2](https://kdl.dev). By default it is called `metapkg.kdl`; use `--file` to choose another name. Comments start with `//`.

Packages are installed in the order they appear in the file.

## Package entries

A package is a line with the package name, optionally followed by a package manager in quotes:

```
bat
fzf "dnf"
starship "snap"
```

The package name is whatever the package manager calls it. The package manager name is always a quoted string.

### Package managers

| Name      | Install command                      | Installed check        |
| --------- | ------------------------------------ | ---------------------- |
| `apt`     | `sudo apt-get install -y <name>`     | `dpkg-query`           |
| `dnf`     | `sudo dnf install -y <name>`         | `dnf list installed`   |
| `snap`    | `sudo snap install <name>`           | `snap list`            |
| `flatpak` | `flatpak install -y <name>`          | `flatpak info <name>`  |
| `brew`    | `brew install --formula\|--cask <name>` | `brew list`         |
| `script`  | the commands in a `scriptdef`        | `which <name>`         |

Notes on individual managers:

- **snap**: if snap reports that a package needs classic confinement, metapkg retries with `--classic`.
- **flatpak**: if flatpak cannot choose a remote, metapkg retries with the `flathub` remote. Flatpak names look like `org.kde.filelight`. A remote can be written after the manager (`md.obsidian.Obsidian "flatpak" "flathub"`), but it is not used yet: the value is currently ignored.
- **brew**: a name can be a formula or a cask. metapkg checks both, so `iterm2` and `bat` need no extra syntax. If a name is both a formula and a cask, or is neither, metapkg stops with an error.
- **script**: see below.

### The default package manager

A package with no manager name uses the default for the operating system, read from `/etc/os-release`:

| Operating system                         | Default |
| ---------------------------------------- | ------- |
| Debian, Ubuntu, and derivatives          | `apt`   |
| Fedora, RHEL, and derivatives            | `dnf`   |
| macOS                                    | `brew`  |

If the system is not recognized, the default is `dnf`. `snap` is never a default; always write `"snap"` explicitly.

You can mix entries freely. The same file can use the default for most packages and name a manager for the rest.

## Script packages

Use a script package for software that has no package in your package manager. A script package needs two parts: an entry marked `"script"`, and a `scriptdef` block with the same name.

```
vivaldi "script"

scriptdef "vivaldi" {
  "sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo"
  "sudo dnf install -y vivaldi-stable"
}
```

Each string in the block is one shell command. The commands are written to a temporary file and run in order with your `$SHELL` (or `/bin/sh` if `$SHELL` is not set).

A command can span several lines. Use a KDL raw string for this, which keeps quotes and line breaks as written:

```
scriptdef "gcloud" {
  #"""
sudo tee -a /etc/yum.repos.d/google-cloud-sdk.repo << EOM
[google-cloud-cli]
name=Google Cloud CLI
baseurl=https://packages.cloud.google.com/yum/repos/cloud-sdk-el9-x86_64
enabled=1
gpgcheck=1
EOM
"""#
  "sudo dnf install -y google-cloud-cli"
}
```

metapkg decides whether a script package is installed by running `which <name>`. That means the package name must be the name of the executable the script installs. If it is already on your `PATH`, the script is skipped.

## A complete example

```
// Command line tools, from the system's default package manager
bat
curl
fzf

// From the snap store
starship "snap"

// A flatpak application
org.kde.filelight "flatpak"

// Installed with a script because there is no package
vivaldi "script"

scriptdef "vivaldi" {
  "sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo"
}
```

More examples are in [`test-files/sample_config.kdl`](../test-files/sample_config.kdl) and [`test-files/sample_essentials.kdl`](../test-files/sample_essentials.kdl).

## Not supported yet

These appear in the sample files as comments, or in the design notes, but do nothing today:

- Groups of packages (`group "name" { ... }`)
- Running a script from a separate file
- Naming the executable a script package installs
- A configuration file for choosing the default package manager
- Removing or updating packages (`metapkg outdated` is a placeholder)

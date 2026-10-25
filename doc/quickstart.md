# Quickstart

This page takes you from nothing to a working install. The [config reference](config-reference.md) describes every part of the file format.

## 1. Get metapkg

Download the latest release for your system (Linux or macOS, `amd64` or `arm64`) with the helper script:

```
curl -fsSL https://raw.githubusercontent.com/kenliu/metapkg/main/get-metapkg.sh | bash
```

This puts a `metapkg` binary in the current directory and checks it against the release checksums. Pass a directory to install somewhere else, for example `curl -fsSL https://raw.githubusercontent.com/kenliu/metapkg/main/get-metapkg.sh | bash -s -- ~/.local/bin`. Set `METAPKG_VERSION=v0.1.0` to get a specific version instead of the latest.

You can also download a binary from the [Releases](https://github.com/kenliu/metapkg/releases) page, or build from source (see the README).

Check that it runs:

```
./metapkg version
```

## 2. Write a metapkg file

Create `metapkg.kdl` in the directory you will run metapkg from. Each line is a package:

```
// Installed with the default package manager for this system
bat
curl
htop

// Installed with a specific package manager
starship "snap"
```

## 3. Preview, then install

Use `--dry-run` first. It checks what is already installed and prints what metapkg would run, without changing anything:

```
./metapkg --dry-run install
```

```
Dry run: no changes will be made
================================
[skip] curl (apt): already installed
[install] bat (apt): would run: sudo apt-get install -y bat
[install] htop (apt): would run: sudo apt-get install -y htop
[install] starship (snap): would run: sudo snap install starship
    may retry with: sudo snap install --classic starship (if snap reports the snap requires classic confinement)

Summary: 3 to install, 1 already installed
Dry run: nothing was changed
```

When the plan looks right, run the install:

```
./metapkg install
```

metapkg goes through the file in order. Packages that are already installed are skipped. If an install fails, metapkg stops with an error and does not try the remaining packages. Run it again after fixing the problem; the packages that succeeded are skipped.

## Using a different file

By default metapkg reads `metapkg.kdl` from the current directory. Use `--file` (or `-f`) to point at another file:

```
./metapkg -f ~/setups/laptop.kdl install
./metapkg install --file ~/setups/laptop.kdl
```

## Commands and options

| Command   | What it does                                               |
| --------- | ---------------------------------------------------------- |
| `install` | Installs the packages in the file that are not installed   |
| `version` | Prints the metapkg version (`dev` for a source build)      |
| `outdated`| Not implemented yet                                        |

| Option            | What it does                                              |
| ----------------- | --------------------------------------------------------- |
| `-f`, `--file`    | The metapkg file to use (default `metapkg.kdl`)           |
| `--dry-run`       | Show what `install` would do without running any installs |
| `--help`          | Show the built-in help                                    |

Options can come before or after the command name.

## Things to know

- `install` runs real commands such as `sudo apt-get install -y`, `sudo dnf install -y`, and `sudo snap install`. Use `--dry-run` when you are not sure what a file will do.
- metapkg only installs. It does not remove packages, update them, or install dependencies beyond what the package manager does itself.
- The package manager must already be present on the machine. A file that lists a `snap` package fails on a system without `snap`.

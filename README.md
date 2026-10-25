# 📦 metapkg

A lightweight, local package management tool for managing packages across different package managers.

## 🤔 What is metapkg?

`metapkg` is a straightforward command line tool designed to help users manage the installation of packages on their local machines across multiple package managers such as `apt`, `dnf`, and `snap`. It provides a simple way to install and track packages with a simple configuration file without the complexity of full-scale package managers such as Ansible or Nix.

## ✨ Features

- 💻 Local package installation
- 🤝 Compatibility with existing system package managers
- 📘 Simple configuration file
- 🧪 Handle scripted installations using system shell commands

## ⚠️ Limitations

To keep things simple, metapkg does not:

- 🚫 Automatically install package dependencies
- 🌐 Manage packages on remote machines
- 🔄 Install package updates directly (use your system package manager for this)
- 🧪 Handle complex installation recipes

## 🚀 Getting Started

Released versions of metapkg for Linux and macOS (`amd64` and `arm64`) can be downloaded from the [Releases](https://github.com/kenliu/metapkg/releases) page. Download the binary for your system (for example `metapkg_linux_amd64`), optionally check it against the `SHA256SUMS` file with `sha256sum -c SHA256SUMS --ignore-missing`, and run `chmod +x` on it.

Prebuilt binaries for Linux and macOS (`amd64` and `arm64`) are also built on every push to `main`. Download them from the artifacts of the latest successful [Build binaries](https://github.com/kenliu/metapkg/actions/workflows/build-binaries.yml) workflow run (you need to be signed in to GitHub to download artifacts). Each binary comes with a `.sha256` checksum file. Artifacts are kept for 30 days, and you need to run `chmod +x` on the binary after unzipping it.

To build metapkg yourself instead, follow these steps to download the Go source and build it:

1. Ensure you have Go installed on your system (version 1.27 or later is recommended).

2. Clone the repository:
   ```
   git clone https://github.com/kenliu/metapkg.git
   cd metapkg
   ```

3. Build the project:
   ```
   go build -o bin ./...
   ```

   This command will compile all packages in the project and its subdirectories, creating an executable named `metapkg` in the `bin` directory.

4. (Optional) Add the bin directory to your PATH for easier access:
   ```
   export PATH=$PATH:$(pwd)/bin
   ```

   You can add this line to your shell configuration file (e.g., `.bashrc`, `.zshrc`) to make it permanent.

Now you can run metapkg by executing `./bin/metapkg` from the project directory, or simply `metapkg` if you've added it to your PATH.

## 📘 Usage

To get started, create a metapkg file called `metapkg.kdl` in the same directory as the `metapkg` executable. You can use this simple example metapkg file to get started:

```
cowsay
sl
htop
```

Then, run the following command to install the packages:

```
./bin/metapkg install
```

This will install the `cowsay`, `sl`, and `htop` packages using the default package manager for your system. The default is chosen from the operating system:

| Operating system                              | Default package manager |
| --------------------------------------------- | ----------------------- |
| Ubuntu, Debian, and derivatives (Mint, Pop!_OS) | `apt`                 |
| Fedora, RHEL, and derivatives (Rocky, Alma)   | `dnf`                   |
| macOS                                         | `brew`                  |

Currently, metapkg supports the following package managers:

- `apt`
- `dnf`
- `flatpak`
- `brew`
- `snap`

If you want to use a different package manager to install a package, you can specify the package manager and any additional arguments as follows:

```
cowsay
sl
htop
org.kde.filelight "flatpak" "flathub"
```

This will additionally install the `org.kde.filelight` package using the `flatpak` package manager with the `flathub` remote.

This is just a simple example. You see a more complete example in the [sample_config.kdl](test-files/sample_config.kdl) file.

### Choosing a metapkg file

By default metapkg reads `metapkg.kdl` from the current directory. To use a different file, pass `--file` (or `-f`):

```
./bin/metapkg --file sample_essentials.kdl install
./bin/metapkg -f ~/setups/laptop.kdl --dry-run install
```

`--file` is a global option, so it works with every command (`install`, `outdated`) and can be written before or after the command name: `metapkg install --file sample_essentials.kdl` is equivalent.

### Dry run

To see what `install` would do without changing anything, pass `--dry-run`. Only read-only checks are run (for example `dpkg-query` or `snap list`); no install command is executed.

```
./bin/metapkg --dry-run install
```

The output lists each package in file order, followed by a summary:

```
Dry run: no changes will be made
================================
[skip] curl (apt): already installed
[install] cowsay (apt): would run: sudo apt-get install -y cowsay
[install] starship (snap): would run: sudo snap install starship
    may retry with: sudo snap install --classic starship (if snap reports the snap requires classic confinement)
[install] vivaldi (script): would run:
    sudo dnf config-manager --add-repo https://repo.vivaldi.com/archive/vivaldi-fedora.repo

Summary: 3 to install, 1 already installed
Dry run: nothing was changed
```

## 🏷️ Releases

To cut a release, a maintainer pushes a version tag from `main`:

```
git checkout main
git pull
git tag v0.2.0
git push origin v0.2.0
```

The [Release](https://github.com/kenliu/metapkg/actions/workflows/release.yml) workflow then tests and builds that commit and publishes a GitHub Release with the binaries and a `SHA256SUMS` file. Tags must look like `v1.2.3` and point to a commit on `main`; tags with a suffix such as `v0.2.0-rc.1` are published as prereleases.

## 🤝 Contributing

Contributions are welcome! If you have ideas for improvements or have found a bug, please feel free to open an issue or submit a pull request.

## 📜 License

`metapkg` is open source software [licensed as MIT](LICENSE). Copyright (c) 2024 Ken Liu.

## 🆘 Support

If you encounter any problems or have questions, please open an issue in this repository.   

## 👏 Acknowledgements

[Recognize contributors and any third-party resources here]

Thank you for checking out `metapkg`. I hope it simplifies your local package management!
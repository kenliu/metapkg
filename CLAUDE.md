# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`metapkg` is a Go CLI that reads a KDL file (`metapkg.kdl`) listing packages and installs the missing ones through the right system package manager (apt, dnf, snap, brew, flatpak, or a shell script). It does not resolve dependencies, update packages, or manage remote machines. The README is the authoritative description of the file format and usage.

## Commands

```
go build -o bin ./...                    # builds bin/metapkg from ./cmd/metapkg (bin/ is gitignored)
go test ./...
go test ./internal/packages -run TestParseApt        # one test
go test ./internal/engine -run TestDryRunInstall -v  # one test with output
```

- `make build` and `make test` work. The `make build-linux/-windows/-mac` targets are stale: they build the module root rather than `./cmd/metapkg`, so use the commands in `.github/workflows/release.yml` for cross-compiling (`CGO_ENABLED=0 GOOS=... GOARCH=... go build -trimpath ./cmd/metapkg`).
- `go vet ./...` reports four known warnings (an `append` with no values in `dnf.go`, and kong struct tags in `main.go`). They predate current work and CI does not run vet. Do not fix them as a side effect of an unrelated change.
- Go version is read from `go.mod` everywhere, including CI. Bump it there only.
- The `version` command prints `dev` unless built with `-ldflags "-X main.version=<tag>"`, which only the release workflow does.
- Running `metapkg install` for real executes `sudo apt-get`, `sudo snap install`, and so on. To see what a file would do, use `metapkg --dry-run install --file X`.

## Architecture

The flow is: KDL file -> `packages.LoadMetapackageFile` -> `Metapackage` -> `engine` -> one `PackageState` per package.

- `cmd/metapkg/main.go` is a kong CLI. `install` calls `engine.InstallPackages`, or `engine.DryRunInstall` when `--dry-run` is set.
- `internal/packages` parses the KDL file. Each top-level node is a package, and its first string argument is the package manager name (`curl "apt"`). A node with no argument uses the default package manager. `scriptdef` nodes are collected separately into `Metapackage.Scriptdefs`.
- `internal/engine` maps a package's manager name to a `PackageState` implementation through `packageManagerFactory` and the `newPackageState` variable. Tests replace `newPackageState` with fakes, so the engine is tested without running any real package manager.
- `internal/package_managers/<name>` holds one implementation per manager, all satisfying the `PackageState` interface in `package_managers.go`: `IsInstalled`, `IsOutdated`, `Install`.

### Adding or changing a package manager

Touching one manager means touching several files, and nothing fails to compile if you miss one:

1. A new package under `internal/package_managers/<name>`.
2. An entry in `packageManagerFactory` in `internal/engine/engine.go`.
3. A `parseX` function and a `case` in the `switch` in `LoadMetapackageFile` in `internal/packages/packages.go`.
4. The supported-manager list in `README.md`.

### Dry run contract

Dry run must never call `Install`; it only runs `IsInstalled` and prints a plan. To describe an install without running it, managers implement the optional `InstallCommander` interface, and `Install` must build its `exec.Command` from `InstallCommand`, so the plan and the real command cannot drift apart. Managers whose install is conditional (snap retries with `--classic`, flatpak falls back to the flathub remote) or generated (script) also implement `InstallDescriber`. Keep any new manager consistent with this.

### Default package manager

A package written with no manager (`bat` alone) uses `packages.DefaultPackageManager()`, which reads `ID` and `ID_LIKE` from `/etc/os-release` (apt for Debian/Ubuntu families, dnf for Fedora/RHEL families, brew on macOS) and falls back to `DEFAULT_PACKAGE_MANAGER` (dnf). Snap is intentionally not a default.

### Script packages

`foo "script"` needs a matching `scriptdef "foo" { ... }` in the same file. Installed state is just `which foo`, so the package name must be the executable name. Commands are written to a temp file and run with `$SHELL`, or `/bin/sh` if it is unset.

## Gotchas

- `Package.Arguments` exists on the struct but no parser fills it in. Extra KDL arguments, such as the flatpak remote in `org.kde.filelight "flatpak" "flathub"`, are currently ignored, even though the README and sample file show them.
- `metapkg outdated` is not implemented (`engine.ListOutdatedPackages` returns "not implemented"), though each manager has an `IsOutdated` method. The `--config`, `--debug`, `--verbose`, and `--quiet` flags are declared but not read.
- `doc/` and `design/` describe planned features (a `metapkg.conf` file, a configurable default manager, groups, removing packages). Much of it is not implemented, so trust the code and README over those files.
- `.gitignore` contains the bare pattern `metapkg`, which also matches the `cmd/metapkg/` directory. Existing files there are tracked and edit fine, but a new file under `cmd/metapkg/` is silently ignored; add it with `git add -f`.
- `TestLoadMetapackage` reads `test-files/sample_config.kdl` and asserts exactly 8 packages and 2 scriptdefs, so editing that sample file means updating the test.
- Manager unit tests cover pure helpers (output parsing, `InstallCommand`), not real installs. Snap, flatpak, brew, and dnf behavior cannot be exercised on a machine without those tools, so say so when you could not verify it.

## CI and releases

- `go.yml` builds and tests on pull requests and pushes to `main`. `dependency-review.yml` runs on pull requests.
- `build-binaries.yml` runs on every push to `main` and uploads linux/darwin amd64/arm64 binaries as workflow artifacts (30 days).
- `release.yml` runs when a tag matching `v*.*.*` is pushed. The tag must be strict semver and point at a commit on `main`. It creates a draft release, attaches the binaries and `SHA256SUMS`, then publishes it. Tags with a hyphen (`v0.2.0-rc.1`) become prereleases. Create the tag with `git tag` and push it; do not create the release in the GitHub UI first, because the workflow refuses to modify an already published release.
- `actions/dependency-review-action` is pinned to `v5.0.0` because that action has no floating `v5` tag. The other actions use floating major tags.

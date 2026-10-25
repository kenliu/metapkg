package main

import (
	"fmt"
	"os"

	"github.com/kenliu/metapkg/internal/engine"
	"github.com/kenliu/metapkg/internal/packages"

	"github.com/alecthomas/kong"
)

type Config struct {
	DefaultPackageManager string `kdl:"default_package_manager"`
}

// version is the version reported by the version command. It is set at build
// time by the release workflow with -ldflags "-X main.version=<tag>". Builds
// that do not set it report "dev".
var version = "dev"

var cli struct {
	Config  string `help:"Config file path" default:"metapkg.conf" type:"path"`
	File    string `help:"Metapkg file to use" default:"metapkg.kdl" type:"path"`
	DryRun  bool   `help:"Print what would change without executing any install commands" default:"false"`
	Debug   bool   `help:"Print debug information" default:"false"`
	Verbose bool   `help:"Print verbose information" default:"false"`
	Quiet   bool   `help:"Suppress all output" default:"false"`

	Install  struct{} `cmd help:"Install packages specified in the metapkg file"`
	Outdated struct{} `cmd help:"List outdated packages"`

	Version struct{} `cmd help:"Print the version number of metapkg"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("metapkg"),
		kong.Description("A tool to install packages using different package managers"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	switch ctx.Command() {
	case "install":
		if err := install(cli.File, cli.DryRun); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "outdated":
		if err := outdated(cli.File); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println("metapkg " + version)
	}
}

func install(file string, dryRun bool) error {
	metapackage, err := packages.LoadMetapackageFile(file)
	if err != nil {
		return fmt.Errorf("error loading metapkg file: %w", err)
	}
	if dryRun {
		return engine.DryRunInstall(metapackage, os.Stdout)
	}
	return engine.InstallPackages(metapackage)
}

func outdated(file string) error {
	metapackage, err := packages.LoadMetapackageFile(file)
	if err != nil {
		return fmt.Errorf("error loading metapkg file: %w", err)
	}
	return engine.ListOutdatedPackages(metapackage)
}

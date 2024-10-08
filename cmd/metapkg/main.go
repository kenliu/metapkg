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

var cli struct {
	Config  string `help:"Config file path" default:"metapkg.conf" type:"path"`
	DryRun  bool   `help:"Print commands without executing" default:"false"`
	Debug   bool   `help:"Print debug information" default:"false"`
	Verbose bool   `help:"Print verbose information" default:"false"`
	Quiet   bool   `help:"Suppress all output" default:"false"`

	Install struct {
		File string `help:"Metapkg file to use" default:"metapkg.kdl" type:"path"`
	} `cmd help:"Install packages specified in metapkg.kdl"`

	Outdated struct {
		File string `help:"Metapkg file to use" default:"metapkg.kdl" type:"path"`
	} `cmd help:"List outdated packages"`

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
		if err := install(cli.Install.File); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "outdated":
		if err := outdated(cli.Outdated.File); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println("metapkg v0.1")
	}
}

func install(file string) error {
	metapackage, err := packages.LoadMetapackageFile(file)
	if err != nil {
		return fmt.Errorf("error loading metapkg file: %w", err)
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

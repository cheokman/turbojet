package main

import (
	"os"
	"turbojet/cli"
	"turbojet/config"
	"turbojet/content"
	"turbojet/provider"
)

func main() {
	writer := cli.DefaultWriter()

	ctx := cli.NewCommandContext(writer)

	rootCmd := &cli.Command{
		Name:  "tj",
		Short: "GameSource Cloud CDN Command Line Interface Version " + cli.Version,
		Usage: "tj <product> <operation> [--parameter1 value1 --parameter2 value2 ...]",
	}
	config.AddFlags(rootCmd.Flags())
	provider.AddFlags(rootCmd.Flags())
	content.AddFlags(rootCmd.Flags())

	ctx.EnterCommand(rootCmd)

	rootCmd.AddSubCommand(cli.NewVersionCommand())
	rootCmd.AddSubCommand(config.NewConfigureCommand())
	rootCmd.AddSubCommand(provider.NewProviderCommand())
	rootCmd.AddSubCommand(content.NewContentCommand())
	rootCmd.Execute(ctx, os.Args[1:])
}

package config

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"turbojet/cli"
)

func NewConfigureListCommand() *cli.Command {
	cmd := &cli.Command{
		Name:  "list",
		Short: "list configuration",
		Usage: "list",
		Run: func(ctx *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], ctx)
			}
			doConfigureList(ctx)
			return nil
		},
	}
	return cmd
}

func doConfigureList(ctx *cli.Context) {
	config, err := LoadConfiguration(GetConfigPath()+string(os.PathSeparator)+configFile, ctx.Writer())
	if err != nil {
		cli.Errorf(ctx.Writer(), "load configuration failed %s", err)
	}
	tw := tabwriter.NewWriter(ctx.Writer(), 8, 0, 1, ' ', 0)
	fmt.Fprint(tw, "\nProfile\t| Domain\t| Valid\n")
	fmt.Fprintf(tw, "-------\t| --------\t| -------\n")
	for _, pv := range config.Profiles {
		name := pv.Name
		err := pv.Validate()
		valid := "Valid"
		if err != nil {
			valid = "Invalid"
		}
		var domains string
		if name == config.CurrentProfile {
			name = name + " *"
			domains = pv.ListDomains()
		} else {
			domains = strings.Join(pv.Domains, ",")
		}
		fmt.Fprintf(tw, "%s\t| %s\t| %s\n", name, domains, valid)
	}
	tw.Flush()
}

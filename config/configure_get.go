package config

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"turbojet/cli"
)

func NewConfigureGetCommand() *cli.Command {
	cmd := &cli.Command{
		Name:  "get",
		Short: "print configuration values",
		Usage: "get [profile] [domain] ...",
		Run: func(c *cli.Context, args []string) error {
			doConfigureGet(c, args)
			return nil
		},
	}
	return cmd
}

func doConfigureGet(c *cli.Context, args []string) {
	config, err := LoadConfiguration(GetConfigPath()+string(os.PathSeparator)+configFile, c.Writer())
	if err != nil {
		cli.Errorf(c.Writer(), "load configuration failed %s", err)
	}

	profile := config.GetCurrentProfile(c)

	if pn, ok := ProfileFlag(c.Flags()).GetValue(); ok {
		profile, ok = config.GetProfile(pn)
		if !ok {
			cli.Errorf(c.Writer(), "profile %s not found!", pn)
		}
	}

	if len(args) == 0 && !reflect.DeepEqual(profile, Profile{}) {
		data, err := json.MarshalIndent(profile, "", "\t")
		if err != nil {
			cli.Printf(c.Writer(), "ERROR:"+err.Error())
		}
		cli.Println(c.Writer(), string(data))
	} else {
		for _, arg := range args {
			switch arg {
			case ProfileFlagName:
				cli.Printf(c.Writer(), "profile=%s\n", profile.Name)
			case DomainFlagName:
				cli.Printf(c.Writer(), "domain=%s\n", strings.Join(profile.Domains[:], ","))
			}
		}
	}

	cli.Printf(c.Writer(), "\n")
}

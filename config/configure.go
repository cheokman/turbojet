package config

import (
	"os"
	"turbojet/cli"
)

func NewConfigureCommand() *cli.Command {
	c := &cli.Command{
		Name:  "configure",
		Short: "configure profile and domain",
		Usage: "configure --profile <ProfileName> --domain <DomainName>",
		Run: func(ctx *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], ctx)
			}

			profileName, _ := ProfileFlag(ctx.Flags()).GetValue()
			domainName, _ := DomainFlag(ctx.Flags()).GetValue()

			return doConfigure(ctx, profileName, domainName)
		},
	}

	c.AddSubCommand(NewConfigureGetCommand())
	// c.AddSubCommand(NewConfigureSetCommand())
	c.AddSubCommand(NewConfigureListCommand())
	// c.AddSubCommand(NewConfigureDeleteCommand())

	return c
}

func doConfigure(ctx *cli.Context, profileName string, domainName string) error {
	w := ctx.Writer()

	conf, err := LoadConfiguration(GetConfigPath()+string(os.PathSeparator)+configFile, ctx.Writer())
	if err != nil {
		return err
	}

	if profileName == "" {
		profileName = DefaultConfigProfileName
	}
	cp, ok := conf.GetProfile(profileName)
	if !ok {
		cp = conf.NewProfile(profileName)
	}

	cli.Printf(w, "Configuring profile '%s' using '%s' domain... \n", profileName, domainName)

	if domainName != "" {
		cp.PutDomain(domainName)
	}

	cli.Printf(w, "Saving profile[%s] ...", profileName)
	conf.PutProfile(cp)
	conf.CurrentProfile = cp.Name
	err = SaveConfiguration(conf)

	if err != nil {
		return err
	}
	cli.Printf(w, "Done.\n")

	return nil
}

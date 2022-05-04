package config

import "turbojet/cli"

const (
	ProfileFlagName       = "profile"
	DomainFlagName        = "domain"
	ConfigurePathFlagName = "config-path"
)

func AddFlags(fs *cli.FlagSet) {
	fs.Add(NewProfileFlag())
	fs.Add(NewDomainFlag())
	fs.Add(NewConfigurePathFlag())
}

func NewProfileFlag() *cli.Flag {
	return &cli.Flag{
		Category:     "config",
		Name:         ProfileFlagName,
		Shorthand:    'p',
		DefaultValue: "default",
		Persistent:   true,
		Short:        "use `--Profile <profileName> to select CDN profile",
	}
}

func NewDomainFlag() *cli.Flag {
	return &cli.Flag{
		Category:   "config",
		Name:       DomainFlagName,
		Persistent: true,
		Short:      "use `--domain <domainName> to select domain of CDN",
	}
}

func NewConfigurePathFlag() *cli.Flag {
	return &cli.Flag{
		Category:     "config",
		Name:         ConfigurePathFlagName,
		AssignedMode: cli.AssignedOnce,
		Persistent:   true,
		Short:        "use `--config-path` to specify the configuration file path",
	}
}

func ProfileFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(ProfileFlagName)
}

func DomainFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(DomainFlagName)
}

func ConfigurePathFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(ConfigurePathFlagName)
}

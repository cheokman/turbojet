package provider

import (
	"turbojet/cli"
)

const (
	NameFlagName     = "name"
	DomainFlagName   = "domain"
	CacheFlagName    = "cache"
	HeadLessFlagName = "headless"
)

func AddFlags(fs *cli.FlagSet) {
	fs.Add(NewNameFlag())
	fs.Add(NewCacheFlag())
	fs.Add(NewHeadLessFlag())
}

func NewHeadLessFlag() *cli.Flag {
	return &cli.Flag{
		Category:     "provider",
		Name:         HeadLessFlagName,
		Shorthand:    'h',
		DefaultValue: "true",
		Persistent:   true,
		Short:        "use `--headless <true/false> to enable/disable headless for browser`",
	}
}

func NewNameFlag() *cli.Flag {
	return &cli.Flag{
		Category:     "provider",
		Name:         NameFlagName,
		Shorthand:    'n',
		DefaultValue: "default",
		Persistent:   true,
		Short:        "use `--name <providerName> to select CDN provider`",
	}
}

func NewDomainFlag() *cli.Flag {
	return &cli.Flag{
		Category:     "provider",
		Name:         DomainFlagName,
		AssignedMode: cli.AssignedRepeatable,
		Shorthand:    'd',
		Persistent:   true,
		Short:        "use `--domain <domainName1> <domainName2> to select domain name`",
	}
}

func NewCacheFlag() *cli.Flag {
	return &cli.Flag{
		Category:     "provider",
		Name:         CacheFlagName,
		AssignedMode: cli.AssignedNone,
		Shorthand:    'c',
		Persistent:   true,
		Short:        "use `--cache to enable cache with local file`",
	}
}

func NameFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(NameFlagName)
}

func CacheFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(CacheFlagName)
}

func DomainFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(DomainFlagName)
}

func HeadLessFlag(fs *cli.FlagSet) *cli.Flag {
	return fs.Get(HeadLessFlagName)
}

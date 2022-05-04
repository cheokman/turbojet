package provider

import (
	"turbojet/cli"
	"turbojet/instrument"
)

const (
	defaultCacheFile = "tj.cat"
)

func NewProviderRefreshCommand() *cli.Command {
	return &cli.Command{
		Name:  "refresh",
		Short: "refresh provider information in cache",
		Usage: "refresh --domain [domain] --cache ...",
		Run: func(c *cli.Context, args []string) error {
			doProviderRefresh(c, args)
			return nil
		},
	}
}

func doProviderRefresh(c *cli.Context, args []string) {
	provider, storage, err := LoadProviderWithContext(c)
	w := c.Writer()
	if err != nil {
		cli.Errorf(w, "Loading provider failed %s\n", err)
	}

	var domains []string
	flagDomains := DomainFlag(c.Flags()).GetValues()
	if len(flagDomains) == 0 {
		for _, d := range provider.Domains {
			domains = append(domains, d.Name)
		}
	} else {
		domains = flagDomains
	}

	cli.Printf(w, "Starting refresh domain: %s\n", domains)
	url := provider.InstrURL
	sources, err := instrument.Refresh(c, domains, url)
	if err != nil {
		cli.Errorf(w, "CDN provider refresh error: %s\n", err)
	}

	if _, ok := CacheFlag(c.Flags()).GetValue(); ok {
		err := provider.CacheSource(c, sources, storage)
		if err != nil {
			cli.Errorf(w, "Cache data error: %s\n", err)
		}
	}
	provider.SaveToStorage(storage)
}

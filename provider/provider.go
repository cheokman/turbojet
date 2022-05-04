package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"turbojet/cli"
	"turbojet/config"
	"turbojet/instrument"
)

const (
	DefaultConfigProviderName = "default"
	DefaultInstrURL           = "/yyyy/index.html"
)

type Provider struct {
	Name          string    `json:"name"`
	InstrURL      string    `json:"InstrURL"`
	Domains       []Domain  `json:"domains"`
	DefaultDomain string    `json:"default_domain"`
	RefreshedAt   time.Time `json:"refreshed_at"`
}

func NewProvider() Provider {
	return Provider{
		Name:     DefaultConfigProviderName,
		InstrURL: DefaultInstrURL,
		Domains:  []Domain{},
	}
}

func NewProviderCommand() *cli.Command {
	c := &cli.Command{
		Name:  "provider",
		Short: "provider information and settings",
		Usage: "provider --name <Provider Name>",
		Run: func(c *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], c)
			}

			providerName, _ := NameFlag(c.Flags()).GetValue()
			// domainName, _ := DomainFlag(c.Flags()).GetValue()
			return doProvider(c, providerName)
		},
	}
	c.AddSubCommand(NewProviderRefreshCommand())
	c.AddSubCommand(NewProviderListCommand())
	return c
}

func doProvider(c *cli.Context, providerName string) error {
	var err error
	w := c.Writer()
	profile, err := config.LoadProfileWithContext(c)
	if err != nil {
		return fmt.Errorf("Configuration failed, use `tj configure` to configure it")
	}
	err = profile.Validate()

	storage, err := LoadStorage(GetStoragePath()+string(os.PathSeparator)+storageFile, w)
	if err != nil {
		return err
	}

	var provider Provider
	if providerName != "" {
		p, ok := storage.GetProvider(providerName)
		if !ok {
			cli.Printf(w, "provider '%s' not found", providerName)
			return nil
		}
		provider = p
	} else {
		providerName = profile.Name
		p, ok := storage.GetProvider(providerName)
		if !ok {
			cli.Printf(w, "Creating provider '%s' from profile\n", providerName)
			p = storage.NewProvider(providerName)
		}
		provider = p
	}

	cli.Printf(w, "Refreshing provider[%s] from profile\n", providerName)
	provider.UpdateWithProfile(profile)
	cli.Printf(w, "Saving provider[%s] ...\n", providerName)
	storage.PutProvider(provider)
	err = SaveStorage(storage)

	if err != nil {
		return err
	}
	cli.Printf(w, "Done.\n")
	return nil
}

func LoadProviderWithContext(c *cli.Context) (Provider, Storage, error) {
	var provider Provider
	var storage Storage
	w := c.Writer()
	profile, err := config.LoadProfileWithContext(c)
	if err != nil {
		return provider, storage, fmt.Errorf("Configuration failed, use `tj configure` to configure it")
	}
	err = profile.Validate()

	storage, err = LoadStorage(GetStoragePath()+string(os.PathSeparator)+storageFile, w)
	if err != nil {
		return provider, storage, err
	}

	providerName := profile.Name
	provider, ok := storage.GetProvider(providerName)
	if !ok {
		cli.Printf(w, "Creating provider '%s' from profile\n", providerName)
		provider = storage.NewProvider(providerName)
	}
	cli.Printf(w, "Refreshing provider[%s] from profile\n", providerName)
	provider.UpdateWithProfile(profile)
	cli.Printf(w, "Saving provider[%s] ...\n", providerName)
	storage.PutProvider(provider)
	err = SaveStorage(storage)

	if err != nil {
		return provider, storage, err
	}
	cli.Printf(w, "Done.\n")

	return provider, storage, nil
}

func (p Provider) SaveToStorage(s Storage) error {
	s.PutProvider(p)
	err := SaveStorage(s)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) UpdateWithProfile(profile config.Profile) error {
	p.DefaultDomain = profile.GetDefaultDomain()
	for _, dn := range profile.Domains {
		contained := false
		for _, d := range p.Domains {
			if d.IsEqualWithName(dn) {
				contained = true
				break
			}
		}
		if contained {
			continue
		}
		domain := NewDomain(dn)
		p.Domains = append(p.Domains, domain)
	}
	return nil
}

func (p *Provider) CacheSource(c *cli.Context, sources []map[string]string, s Storage) error {
	for _, sd := range sources {
		if sd["err"] != "" {
			cli.Errorf(c.Writer(), "refresh[%s] source error: %s\n", sd["err"])
			continue
		}
		dn := sd["domain"]
		cli.Printf(c.Writer(), "Writing domain: %s\n", dn)
		index := strings.Index(dn, "//")

		file := dn[index+2:]
		path := GetStoragePath() + string(os.PathSeparator) + file + ".json"
		cli.Printf(c.Writer(), "Writing to file[%s}\n", path)
		err := ioutil.WriteFile(path, []byte(sd["source"]), 0644)
		if err != nil {
			cli.Errorf(c.Writer(), "refresh[%s] save file error: %s\n", err)
		} else {
			domain := p.GetDomainByName(dn)
			domain.SetCacheFile(file + ".json")
		}
	}
	return nil
}

func (p *Provider) LoadSource(c *cli.Context, domain string) (instrument.InstrumentSlice, error) {
	index := strings.Index(domain, "//")
	file := domain[index+2:]
	path := GetStoragePath() + string(os.PathSeparator) + file + ".json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		cli.Errorf(c.Writer(), "loading cache for domain[%s] from file[%s] error: %s\n", domain, file, err)
	}

	instruments, err := instrument.LoadCache(c, content)
	if err != nil {
		return instruments, err
	}
	return instruments, nil
}

func (p *Provider) GetDomainByName(n string) *Domain {
	for i, d := range p.Domains {
		if d.IsEqualWithName(n) {
			return &p.Domains[i]
		}
	}

	return &Domain{Name: n}
}

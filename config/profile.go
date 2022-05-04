package config

import (
	"fmt"
	"os"
	"turbojet/cli"
)

type Profile struct {
	Name          string   `json:"name"`
	Domains       []string `json:"domains"`
	defaultDomain string   `json:"default_domain"`
}

func NewProfile(name string) Profile {
	return Profile{
		Name: name,
	}
}

func (p *Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name can't be empty")
	}

	return nil
}

func (p *Profile) PutDomain(domain string) {
	for _, d := range p.Domains {
		if d == domain {
			return
		}
	}
	p.Domains = append(p.Domains, domain)
}

func (p *Profile) SetDefaultDomain(domain string) bool {
	for _, d := range p.Domains {
		if d == domain {
			p.defaultDomain = domain
			return true
		}
	}
	return false
}

func (p *Profile) GetDefaultDomain() string {
	if p.defaultDomain == "" {
		return p.Domains[0]
	}
	return p.defaultDomain
}

func (p *Profile) ListDomains() string {
	var defaultDom, domainStr string
	if p.defaultDomain == "" {
		defaultDom = p.Domains[0]
	} else {
		defaultDom = p.defaultDomain
	}

	domainStr = defaultDom + " *"

	for _, d := range p.Domains {
		if d != defaultDom {
			domainStr = domainStr + "," + d
		}
	}
	return domainStr
}

func (p *Profile) OverwriteWithFlags(ctx *cli.Context) {
	if envDM := os.Getenv("TJ_DOMAIN"); envDM != "" {
		p.PutDomain(envDM)
	}
}

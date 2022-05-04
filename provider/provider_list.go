package provider

import (
	"fmt"
	"text/tabwriter"
	"turbojet/cli"
)

func NewProviderListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Short: "list provider information",
		Usage: "list --provider",
		Run: func(c *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], c)
			}
			doProviderList(c)
			return nil
		},
	}
}

func doProviderList(c *cli.Context) {
	_, storage, err := LoadProviderWithContext(c)

	w := c.Writer()

	if err != nil {
		cli.Errorf(w, "Loading provider failed %s\n", err)
	}

	tw := tabwriter.NewWriter(w, 8, 0, 1, ' ', 0)
	fmt.Fprint(tw, "\nProvider\t| Domain\t| Regions\t| Nodes\t|Valid \n")
	fmt.Fprint(tw, "--------\t| --------\t| -------\t| --------\t| --------\n")
	providers := storage.GetProviders()
	for _, pv := range providers {
		domains := pv.Domains
		valid := "Valid"
		for _, dm := range domains {
			dmN := dm.Name
			instrS, err := pv.LoadSource(c, dmN)
			if err != nil {
				valid = "Invalid"
				fmt.Fprintf(tw, "%s\t| %s\t| %s\t| %s\t| %s\n", pv.Name, dmN, " x", " x", valid)
				continue
			}
			fmt.Fprintf(tw, "%s\t| %s\t| %d\t| %d\t| %s\n", pv.Name, dmN, len(instrS.GetProvinces()), len(instrS.GetIPs()), valid)
		}
	}
	tw.Flush()
}

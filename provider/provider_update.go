package provider

import (
	"fmt"
	"turbojet/cli"
)

func NewProviderUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Short: "update provider information",
		Usage: "update --from-cache ...",
		Run: func(c *cli.Context, args []string) error {
			doProviderUpdate(c, args)
			return nil
		},
	}
}

func doProviderUpdate(c *cli.Context, args []string) {
	provider, storage, err := LoadProviderWithContext(c)
	w := c.Writer()

	if err != nil {
		cli.Errorf(w, "Loading provider failed %s\n", err)
	}

	fmt.Println(provider, storage)
}

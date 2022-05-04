package cli

import "strings"

var (
	Version = "0.0.1"
)

func GetVersion() string {
	return strings.Replace(Version, " ", "-", -1)
}

func NewVersionCommand() *Command {
	return &Command{
		Name:  "version",
		Short: "print current version",
		Run: func(ctx *Context, args []string) error {
			Printf(ctx.Writer(), "%s\n", Version)
			return nil
		},
	}
}

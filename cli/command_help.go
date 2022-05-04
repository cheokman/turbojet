package cli

import (
	"fmt"
	"text/tabwriter"
)

func (c *Command) PrintHead(ctx *Context) {
	Printf(ctx.Writer(), "%s\n", c.Short)

}

func (c *Command) PrintUsage(ctx *Context) {
	if c.Usage != "" {
		Printf(ctx.Writer(), "\nUsage:\n  %s\n", c.GetUsageWithParent())
	} else {
		c.PrintSubCommands(ctx)
	}
}

func (c *Command) PrintSample(ctx *Context) {
	if c.Sample != "" {
		Printf(ctx.Writer(), "\nSample:\n  %s\n", c.Sample)
	}
}

func (c *Command) PrintSubCommands(ctx *Context) {
	if len(c.subCommands) > 0 {
		Printf(ctx.Writer(), "\nCommands:\n")
		w := tabwriter.NewWriter(ctx.Writer(), 8, 0, 1, ' ', 0)
		for _, cmd := range c.subCommands {
			fmt.Fprintf(w, "  %s\t%s\n", cmd.Name, cmd.Short)
		}
		w.Flush()
	}
}

func (c *Command) PrintTail(ctx *Context) {
	Printf(ctx.Writer(), "\nUse `%s --help` for more information.\n", c.Name)
}

package cli

import "fmt"

type Command struct {
	Name string

	Short string

	Usage string

	Sample string

	Run func(ctx *Context, args []string) error

	parent *Command

	EnableUnknownFlag bool

	subCommands []*Command
	flags       *FlagSet
}

func (c *Command) AddSubCommand(cmd *Command) {
	cmd.parent = c
	c.subCommands = append(c.subCommands, cmd)
}

func (c *Command) Flags() *FlagSet {
	if c.flags == nil {
		c.flags = NewFlagSet()
	}
	return c.flags
}

func (c *Command) Execute(ctx *Context, args []string) {
	err := c.executeInner(ctx, args)

	if err != nil {
		c.processError(ctx, err)
	}
}

func (c *Command) GetSubCommand(s string) *Command {
	for _, cmd := range c.subCommands {
		if cmd.Name == s {
			return cmd
		}
	}
	return nil
}

func (c *Command) GetUsageWithParent() string {
	usage := c.Usage
	for p := c.parent; p != nil; p = p.parent {
		usage = p.Name + " " + usage
	}

	return usage
}

func (c *Command) executeInner(ctx *Context, args []string) error {

	parser := NewParser(args, ctx)

	nextArg, _, err := parser.ReadNextArg()
	if err != nil {
		return err
	}

	if nextArg == "help" {
		ctx.help = true
		return c.executeInner(ctx, parser.GetRemains())
	}

	if nextArg != "" {
		subCommand := c.GetSubCommand(nextArg)
		if subCommand != nil {
			ctx.EnterCommand(subCommand)
			return subCommand.executeInner(ctx, parser.GetRemains())
		}

		if c.Run == nil {
			return NewInvalidCommandError(nextArg, ctx)
		}
	}
	remainArgs, err := parser.ReadAll()
	if err != nil {
		return fmt.Errorf("parse failed %s", err)
	}

	err = ctx.CheckFlags()
	if err != nil {
		return err
	}
	if HelpFlag(ctx.Flags()).IsAssigned() {
		ctx.help = true
	}

	callArgs := make([]string, 0)
	if nextArg != "" {
		callArgs = append(callArgs, nextArg)
	}
	for _, s := range remainArgs {
		if s != "help" {
			callArgs = append(callArgs, s)
		} else {
			ctx.help = true
		}
	}

	if ctx.help {
		c.executeHelp(ctx, callArgs)
		return nil
	} else if c.Run == nil {
		c.executeHelp(ctx, callArgs)
		return nil
	} else {
		return c.Run(ctx, callArgs)
	}
}

func (c *Command) processError(ctx *Context, err error) {
	Printf(ctx.Writer(), "ERROR: %s\n", err.Error())
	Exit(1)
}

func (c *Command) executeHelp(ctx *Context, args []string) {
	c.PrintHead(ctx)
	c.PrintUsage(ctx)
	c.PrintSubCommands(ctx)
	c.PrintTail(ctx)
}

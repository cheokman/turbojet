package cli

import (
	"fmt"
	"io"
)

func NewHelpFlag() *Flag {
	return &Flag{
		Name:         "help",
		Short:        "print help",
		AssignedMode: AssignedNone,
	}
}

type Context struct {
	help         bool
	flags        *FlagSet
	unknownFlags *FlagSet
	command      *Command
	writer       io.Writer
}

func NewCommandContext(w io.Writer) *Context {
	return &Context{
		flags:        NewFlagSet(),
		unknownFlags: nil,
		writer:       w,
	}
}

func (ctx *Context) IsHelp() bool {
	return ctx.help
}

func (ctx *Context) Writer() io.Writer {
	return ctx.writer
}

func (ctx *Context) Flags() *FlagSet {
	return ctx.flags
}

func (ctx *Context) UnknownFlags() *FlagSet {
	return ctx.unknownFlags
}

func HelpFlag(fs *FlagSet) *Flag {
	return fs.Get("help")
}

func (ctx *Context) EnterCommand(cmd *Command) {
	ctx.command = cmd
	if !cmd.EnableUnknownFlag {
		ctx.unknownFlags = nil
	} else if ctx.unknownFlags == nil {
		ctx.unknownFlags = NewFlagSet()
	}

	ctx.flags = cmd.flags.mergeWith(ctx.flags, func(f *Flag) bool {
		return f.Persistent
	})
	ctx.flags.Add(NewHelpFlag())
}

func (ctx *Context) CheckFlags() error {
	for _, f := range ctx.flags.Flags() {
		if !f.IsAssigned() {
			if f.Required {
				return fmt.Errorf("missing flag --%s", f.Name)
			}
		} else {
			if err := f.checkFields(); err != nil {
				return err
			}
			if len(f.ExcludeWith) > 0 {
				for _, es := range f.ExcludeWith {
					if _, ok := ctx.flags.GetValue(es); ok {
						return fmt.Errorf("flag --%s is exclusive with --%s", f.Name, es)
					}
				}
			}
		}
	}
	return nil
}

func (ctx *Context) detectFlag(name string) (*Flag, error) {
	flag := ctx.flags.Get(name)
	if flag != nil {
		return flag, nil
	} else if ctx.unknownFlags != nil {
		return ctx.unknownFlags.AddByName(name)
	} else {
		return nil, NewInvalidFlagError(name, ctx)
	}
}

func (ctx *Context) detectFlagByShorthand(ch rune) (*Flag, error) {
	flag := ctx.flags.GetByShorthand(ch)
	if flag != nil {
		return flag, nil
	}
	return nil, fmt.Errorf("unknown flag -%s", string(ch))
}

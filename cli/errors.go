package cli

import "fmt"

type InvalidCommandError struct {
	Name string
	ctx  *Context
}

func NewInvalidCommandError(name string, ctx *Context) error {
	return &InvalidCommandError{
		Name: name,
		ctx:  ctx,
	}
}

func (e *InvalidCommandError) Error() string {
	return fmt.Sprintf("'%s' is not a valid command ", e.Name)
}

type InvalidFlagError struct {
	Flag string
	ctx  *Context
}

func NewInvalidFlagError(name string, ctx *Context) error {
	return &InvalidFlagError{
		Flag: name,
		ctx:  ctx,
	}
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf("invalid flag %s", e.Flag)
}

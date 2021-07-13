package arguments

import (
	"fmt"
	"strings"
)

type Arguments interface {
	// Command gets the first argument from the command line, only if it is NOT a flag.
	// If the cmd line begins with a flag argument, command returns empty
	Command() string

	// Argument gets a named argument from the command line.
	// The name is case insensitive and should match an argument in the command line
	// The resulting argument contains that name, as the given name, with all following argument, up to but excluding tany following flags.
	Argument(name string) *Argument

	// Flags gets all the Arguments with their parameters, with names begining with a '-'
	Flags() []*Argument

	// Remove removes the given argument from the command line.
	// The argument and all its parameters are removed from any further queirs to the Arguments.
	Remove(arg *Argument) error

	// CommandLine gets the complete command line
	CommandLine() []string
}

type Argument struct {
	Name       string
	Parameters []string
}

func (a arguments) CommandLine() []string {
	return a.cmdline
}

type arguments struct {
	cmdline []string
}

func (a arguments) IsEmpty() bool {
	return len(a.cmdline) == 0
}

func (a arguments) Command() string {
	if a.IsEmpty() || strings.HasPrefix(a.cmdline[0], "-") {
		// no command, all flags or empty
		return ""
	}
	return a.cmdline[0]
}

func (a arguments) Flags() []*Argument {
	var flags []*Argument
	for i, cmd := range a.cmdline {
		if !strings.HasPrefix(cmd, "-") {
			continue
		}
		arg := a.newArg(cmd, i)
		i += len(arg.Parameters)
		flags = append(flags, arg)
	}
	return flags
}

// Argument locates the named argument (case insensitive) and gathers any parameters following it.
// parameters begin with the first arg found after the matching name and any following that, upto the end of the command line or
// a '-'flag arg is encountered.
// e.g. cmd dothisthing -flag1 hello world -flag2 false
// cmd has a single string parameter, -flag1 has 2 string parameters, -flag2 has a single bool param.
func (a arguments) Argument(name string) *Argument {
	for i, arg := range a.cmdline {
		if strings.EqualFold(arg, name) {
			return a.newArg(name, i)
		}
	}
	return nil
}

func (a *arguments) Remove(arg *Argument) error {
	i := a.Position(arg)
	if i < 0 {
		return fmt.Errorf("unknown argument %s", arg.Name)
	}
	if i+len(arg.Parameters) >= len(a.cmdline) {
		return fmt.Errorf("invaid arg %s. parameters not found in command line", arg.Name)
	}
	a.cmdline = append(a.cmdline[:i], a.cmdline[i+len(arg.Parameters)+1:]...)
	return nil
}

// parameters colelcts all the arguments following the given position, if any.
// parameters are all arguments following which do NOT start with a '-' flag indicator.
func (a arguments) parameters(position int) []string {
	var params []string
	for i := position + 1; i < len(a.cmdline); i++ {
		// Stop gathering parameters at the next flag or end of cmdline
		if strings.HasPrefix(a.cmdline[i], "-") {
			break
		}
		params = append(params, a.cmdline[i])
	}
	return params
}

func (a arguments) newArg(name string, position int) *Argument {
	return &Argument{
		Name:       name,
		Parameters: a.parameters(position),
	}
}

func (a *arguments) Position(arg *Argument) int {
	for i, a := range a.cmdline {
		if a == arg.Name {
			return i
		}
	}
	return -1
}

func NewArguments(args []string) Arguments {
	return &arguments{cmdline: args}
}

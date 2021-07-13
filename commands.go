package commandgo

import (
	"commandgo/arguments"
	"commandgo/functions"
	"commandgo/help"
	"commandgo/values"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var ErrorNoCommandFound = errors.New("no command found")
var ErrorCommandNotKnown = errors.New("command not known")

// Commands is the main command mapping, mapping command line arguments to variables and functions
// keys must have values of either:
// - A pointer to a global variable or field in an instance of a structure.
// - A function or method on an instance of a structure.
// - Another Commands map.  Sub maps are invoked when the key command from the parent map is called.
// A key may be an empty string, indicating it as the default mapping for that map.
// i.e. if the first command arg is unknown, it is treated as a parameter when invoking the default mapping
// Assignments are only applied at each map level. i.e. top level mappings are assigned first, then any sub map assignments afterwards.
// Only when all the assignments have been set is the final func/method mapping invoked.
type Commands map[string]interface{}

type flagMap map[string]*arguments.Argument

// RunArgs executes this commands using the os.Args array as the arguments to parse.
// Same as calling Run(os.Args[1:])
func (c Commands) RunArgs() ([]interface{}, error) {
	return c.Run(os.Args[1:]...)
}

// Run executes this commands using the given argument array
// All arguments mapped to assignments (variables or fields) are extracted from the given array and applied.
// All remaining arguments are used to call a command, the first being the command and any following are used as parameters for that call.
func (c Commands) Run(args ...string) ([]interface{}, error) {
	// help flags are added to indicate if help requested.
	// These prevents all other flags and commands being invoked.
	if _, ok := c[help.HelpFlagShort]; !ok {
		c[help.HelpFlagShort] = &help.HelpRequested
	}
	if _, ok := c[help.HelpFlagFull]; !ok {
		c[help.HelpFlagFull] = &help.HelpRequested
	}

	var result []interface{}

	// collect any flags from cmdline that are mapped in this map (removes them from args)
	cargs := arguments.NewArguments(args)
	flags, err := c.matchFlags(cargs)
	if err != nil {
		return nil, err
	}

	// Invoke all the flags before invoking the command
	v, err := c.invokeFlags(flags)
	if err != nil {
		return nil, err
	}
	result = append(result, v...)

	// Establish the command key, if any
	ca := cargs.Command() // may be empty
	k, ok := c.findKey(ca)
	if !ok {
		// not known, check if default key available
		k, ok = c.findKey("")
	} else {
		if ca != "" {
			if err := cargs.Remove(&arguments.Argument{Name: ca}); err != nil {
				return nil, err
			}
		}
	}

	if help.HelpRequested {
		return help.ShowHelp(k, args...), nil
	}
	if !ok {
		if ca != "" {
			return nil, ErrorCommandNotKnown
		}
		return nil, ErrorNoCommandFound
	}

	cmd := c[k]
	// ensure all flags have been consumed if not jumping into a submap
	if !c.isSubmap(cmd) {
		f := cargs.Flags()
		if len(f) > 0 {
			names := make([]string, len(f))
			for i, fn := range f {
				names[i] = fn.Name
			}
			return nil, fmt.Errorf("unexpected flag found, %s", strings.Join(names, ","))
		}
	}
	v, err = c.invokeCommand(cmd, cargs.CommandLine())
	if err != nil {
		return nil, err
	}
	return append(result, v...), nil
}

// invokeCommand executes the given command, using the given arguments.
// returns any output from the command or an error
func (c Commands) invokeCommand(cmd interface{}, args []string) ([]interface{}, error) {
	if c.isSubmap(cmd) {
		return (cmd.(Commands)).Run(args...)
	}

	if c.isAssignment(cmd) {
		var a string
		if len(args) > 0 {
			a = args[0]
		}
		return nil, values.SetValue(cmd, a)
	}

	if functions.IsFunc(cmd) {
		return functions.CallFunc(cmd, args...)
	}
	return nil, fmt.Errorf("command is mapped to an unknown type %T", cmd)
}

// invokeFlags executes the command of all the given flags.
// Assignments (var/field pointers) are executed first, followed by any remaining func/method mappings.
// returns any return values from the func mappings or an error
func (c Commands) invokeFlags(flags flagMap) ([]interface{}, error) {
	// Check for help first to prevent others being invokes
	hc := flags.HelpCommand()
	if hc != nil {
		return c.invokeCommand(hc, nil)
	}

	funcM := map[string]*arguments.Argument{}
	// perform the assignments first
	for k, arg := range flags {
		cmd := c[k]
		if !c.isAssignment(cmd) {
			funcM[k] = arg
			continue
		}
		_, err := c.invokeCommand(cmd, arg.Parameters)
		if err != nil {
			return nil, err
		}
	}
	// perform any remaining flag functions,
	var result []interface{}
	for k, arg := range funcM {
		iv, err := c.invokeCommand(c[k], arg.Parameters)
		if err != nil {
			return nil, err
		}
		result = append(result, iv)
	}
	return result, nil
}

// matches any flags found in the given arguments, with mapped flags in this Commands.
// Any matched arguments are removed from the given args and copied to the resulting map.
// returns a map keyed with the 'real' (not the command line arg) keys of this commands, mapping to the matching Argument
func (c Commands) matchFlags(args arguments.Arguments) (flagMap, error) {
	m := flagMap{}
	flags := args.Flags()
	for _, arg := range flags {
		k, ok := c.findKey(arg.Name)
		if !ok {
			continue
		}
		var err error
		arg.Parameters, err = c.trimParameters(c[k], arg.Parameters)
		if err != nil {
			return nil, err
		}
		m[k] = arg
		if err := args.Remove(arg); err != nil {
			log.Fatalln(err)
		}
	}
	return m, nil
}

// findKey finds a key from an argumenet in a case insensitive search
func (c Commands) findKey(arg string) (string, bool) {
	for k := range c {
		if strings.EqualFold(k, arg) {
			return k, true
		}
	}
	return "", false
}

func (c Commands) trimParameters(cmd interface{}, parameters []string) ([]string, error) {
	if !c.isAssignment(cmd) {
		return parameters, nil
	}

	if len(parameters) > 1 {
		parameters = parameters[0:1]
	}

	// Special case for bools, which have optional parameter
	if values.IsKind(cmd, reflect.Bool) {
		// see if following parameter is, in fact a bool otherwise don't use it.
		if len(parameters) > 0 {
			if _, err := strconv.ParseBool(parameters[0]); err != nil {
				parameters = parameters[:0]
			}
		}
	} else {
		if len(parameters) != 1 {
			return nil, fmt.Errorf("no parameter value found for flag")
		}
	}
	return parameters, nil

}

func (c Commands) isAssignment(cmd interface{}) bool {
	return reflect.TypeOf(cmd).Kind() == reflect.Ptr && !functions.IsFunc(cmd)
}

func (c Commands) isSubmap(cmd interface{}) bool {
	_, ok := cmd.(Commands)
	return ok
}

func (m flagMap) HelpCommand() interface{} {
	cmd, ok := m[help.HelpFlagShort]
	if ok {
		return cmd
	}
	cmd, ok = m[help.HelpFlagFull]
	if ok {
		return cmd
	}
	return nil
}

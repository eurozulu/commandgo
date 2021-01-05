package mainline

import (
	"fmt"
	"github.com/eurozulu/mainline/reflection"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Commands maps one or more 'command' strings to methods on a mapped struct.
type Commands map[string]interface{}

// Run attempts to call the mapped method, using the first given argument as the key to the command map.
// If the given key is found, the remaining arguments are parsed into flags and parameters before the mapped method is called.
func (cmds Commands) Run(args ...string) error {
	// strip leading arg if it's program name
	if len(args) > 0 && args[0] == os.Args[0] {
		args = args[1:]
	}
	if len(args) == 0 || args[0] == "help" {
		ShowCommands(cmds, args...)
		return nil
	}

	cmd, err := cmds.findCommand(args[0])
	if err != nil {
		return err
	}

	// call with args, less the initial command
	return cmds.callCommand(cmd, args[1:]...)
}

// callCommand parses the given arguments into flags and parameters for the cmdObject's method, then calls that methed using the parsed data
// Flag values are assigned to mapped Fields in the given command object prior to the call.
func (cmds Commands) callCommand(cmd *command, args ...string) error {
	// Get a value of the struct
	val := reflect.ValueOf(cmd.cmdObject)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("non-struct pointer passed to callCommand")
	}

	// parse args for flags and assign to struct fields
	params, err := parseFlags(val, args)
	if err != nil {
		return err
	}

	// using remaining args, parse into parameters for the cmdObject method.
	inParams, err := parseParameters(*cmd.method, params)
	if err != nil {
		return err
	}

	outVals := val.MethodByName(cmd.method.Name).Call(inParams)

	// check if an error returned
	errInterface := reflect.TypeOf((*error)(nil)).Elem()
	for _, ov := range outVals {
		if (ov.Kind() != reflect.Interface) || ov.IsNil() {
			continue
		}
		if ov.Type().Implements(errInterface) {
			return ov.Interface().(error)
		}
	}
	return nil
}

func (cmds Commands) findCommand(arg string) (*command, error) {
	cm, err := newCommandMap(cmds)
	if err != nil {
		return nil, err
	}
	c, ok := cm[arg]
	if !ok || c == nil {
		return nil, fmt.Errorf("%s is not a known command", arg)
	}
	return c, nil
}

// parseParameters parses the given argument slice of strings into a list of Values of the correct type
// for the given method parameters
func parseParameters(m reflect.Method, args []string) ([]reflect.Value, error) {
	sig, err := reflection.NewSignatureOf(m)
	if err != nil {
		return nil, err
	}

	var vals []reflect.Value
	for i, pt := range sig.ParamTypes {
		var val interface{}
		var err error
		// if last param and variadic, wrap final args into a single array
		if sig.IsVariadic && i == len(sig.ParamTypes)-1 {
			if i < len(args) { // optional params provided
				vps, err := variadicParams(args[i:], pt.Elem())
				if err != nil {
					return nil, err
				}
				vals = append(vals, vps...)
			}
			continue

		} else if i < len(args) {
			val, err = reflection.ValueFromString(args[i], pt)
		} else {
			return nil, fmt.Errorf("%s missing argument %d, requires a %s value", m.Name, i, pt.Name())
		}
		if err != nil {
			return nil, fmt.Errorf("argument %d, %v", i, err)
		}
		vals = append(vals, reflect.ValueOf(val))
	}
	return vals, nil
}

func variadicParams(args []string, t reflect.Type) ([]reflect.Value, error) {
	vals := make([]reflect.Value, len(args))
	for i, arg := range args {
		val, err := reflection.ValueFromString(arg, t)
		if err != nil { // failed to parse as correct type, not a match
			return nil, fmt.Errorf("parameter %v could not be parsed as a %v", arg, t.String())
		}
		vals[i] = reflect.ValueOf(val)
	}
	return vals, nil

}

// parseFlags parses the given arguments for '-' flags, named values, assigning
// any named value to a field of the same name (or tagged as that name) in the given value structure.
func parseFlags(val reflect.Value, args []string) ([]string, error) {
	var unnamed []string
	for i := 0; i < len(args); i++ {
		// collect non flag parameters
		if !strings.HasPrefix(args[i], "-") && args[i] != "-" {
			unnamed = append(unnamed, args[i])
			continue
		}

		// Locate field in struct of the flag name
		arg := strings.TrimLeft(args[i], "-")
		fld := reflection.FindFieldByName(arg, val.Type().Elem(), reflection.FlagTag)
		if fld == nil {
			return nil, fmt.Errorf("--%s is an unknown flag", arg)
		}

		// specical case for bool flags as only one with optional parameter
		if fld.Type.Kind() == reflect.Bool {
			var b = "true"
			// test if following arg exists and is bool
			if i+1 < len(args) {
				_, err := strconv.ParseBool(args[i+1])
				if err == nil {
					b = args[i+1]
					i++
				}
			}
			if err := setFlagValue(b, val.Elem().FieldByName(fld.Name)); err != nil {
				return nil, err
			}
			continue
		}

		// All following ars must have following parameter value
		if i+1 >= len(args) {
			return nil, fmt.Errorf("missing value for flag -%s", arg)
		}
		i++

		if err := setFlagValue(args[i], val.Elem().FieldByName(fld.Name)); err != nil {
			return nil, err
		}
	}
	return unnamed, nil
}

func setFlagValue(v string, fld reflect.Value) error {
	iv, err := reflection.ValueFromString(v, fld.Type())
	if err != nil {
		return fmt.Errorf("invalid flag %s value  %v", fld.Type().String(), err)
	}

	if err := reflection.SetFieldValue(fld, iv); err != nil {
		return err
	}
	return nil
}

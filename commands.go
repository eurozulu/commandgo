package mainline

import (
	"fmt"
	"github.com/eurozulu/mainline/reflection"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Commands maps one or more 'command' strings to methods on a mapped struct.
type Commands map[string]interface{}

// Run attempts to call the mapped method, using the first given argument as the key to the command map.
// If the given key is found, the remaining arguments are parsed into flags and parameters before the mapped method is called.
// returns the return values of the method or an error if the method can not be called or the method returns an error
func (cmds Commands) Run(args ...string) ([]interface{}, error) {
	// strip leading arg if it's program name
	if len(args) > 0 && args[0] == os.Args[0] {
		args = args[1:]
	}
	if len(args) == 0 {
		cmds.ShowCommands()
		return nil, nil
	}

	cmd, err := cmds.findCommand(args[0])
	if err != nil {
		return nil, err
	}

	// call with args, less the initial command
	return cmds.callCommand(cmd, args[1:]...)
}

// ShowCommands lists all the avilable commands and their aliases
func (cmds Commands) ShowCommands() {
	var c []string
	for k := range cmds {
		ks := strings.Split(k, ",")
		s := ks[0]
		if strings.HasPrefix(s, "-") {
			if len(ks) < 2 {
				return
			}
			s = ks[1]
		}
		if len(ks) > 1 {
			s = strings.Join([]string{s, fmt.Sprintf("\t\t(%s)", strings.Join(ks[1:], ", "))}, "")
		}
		c = append(c, s)
	}
	sort.Strings(c)
	for _, cmd := range c {
		fmt.Println(cmd)
	}
}

// callCommand parses the given arguments into flags and parameters for the cmdObject's method, then calls that methed using the parsed data
// Flag values are assigned to mapped Fields in the given command object prior to the call.
func (cmds Commands) callCommand(cmd *command, args ...string) ([]interface{}, error) {
	// Get a value of the struct
	val := reflect.ValueOf(cmd.cmdObject)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("non-struct pointer passed to callCommand")
	}

	// parse args for flags and assign to struct fields
	params, err := parseFlags(val, args)
	if err != nil {
		return nil, err
	}

	// using remaining args, parse into parameters for the cmdObject method.
	inParams, err := parseParameters(*cmd.method, params)
	if err != nil {
		return nil, err
	}

	outVals := val.MethodByName(cmd.method.Name).Call(inParams)

	// flip out values into interface slice and check if an error returned
	var outi []interface{}
	errInterface := reflect.TypeOf((*error)(nil)).Elem()
	for _, ov := range outVals {
		if (ov.Kind() == reflect.Ptr || ov.Kind() == reflect.Interface) && ov.IsNil() {
			continue
		}
		t := ov.Type()
		if t.Kind() == reflect.Interface {
			if t.Implements(errInterface) {
				return nil, ov.Interface().(error)
			}
			outi = append(outi, ov.Elem())
			continue
		}
		outi = append(outi, ov.Interface())
	}
	return outi, nil
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
				val, err = reflection.ValueFromString(strings.Join(args[i:], ","), pt)
			}
		} else if i < len(args) {
			val, err = reflection.ValueFromString(args[i], pt)
		} else {
			return nil, fmt.Errorf("%s missing argument %d, requires a %s value", m.Name, i, pt.Name())
		}
		if err != nil {
			return nil, fmt.Errorf("argument %d, %v", i, err)
		}
		vals = append(vals, reflect.ValueOf(val).Elem())
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

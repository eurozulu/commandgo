package mainline

import (
	"fmt"
	"github.com/eurozulu/mainline/reflection"
	"os"
	"reflect"
	"runtime"
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
	if len(args) == 0 {
		return fmt.Errorf("no command found.  provide a command name or use 'help' to find out more")
	}

	// find the mapped command from the first arg
	cmd := cmds.findCommand(args[0])
	i, ok := cmds[cmd]
	if !ok {
		return fmt.Errorf("'%s' is not a known command", args[0])
	}
	args = args[1:]

	// Get type of the struct and method
	st, err := structFromMethodFunc(i)
	if err != nil {
		return fmt.Errorf("Command Configuration Error:  %v", err)
	}
	md, err := methodFromMethodFunc(i, st)
	if err != nil {
		return err
	}

	// new instance of struct
	ns := reflect.New(st)

	// parse args for flags and assign to struct fields
	params, err := parseFlags(ns, args)
	if err != nil {
		return err
	}

	inParams, err := parseParameters(*md, params)
	if err != nil {
		return err
	}
	outVals := ns.MethodByName(md.Name).Call(inParams)

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

// findCommand looks through the map keys in non case sensative search
// returns the case sensative key if found or empty if not present
func (cmds Commands) findCommand(arg string) string {
	for k := range cmds {
		if strings.EqualFold(k, arg) {
			return k
		}
	}
	return ""
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
	if !sig.IsVariadic && len(vals) < len(args) {
		return nil, fmt.Errorf("too many arguments.  %v expected, found %v", sig.String(), args)
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

// parseFlags parses the given arguments of strings for '-' flags, named values, assigning
// any named value to a field of the same name (or tagged as that name) in the given value structure.
func parseFlags(stc reflect.Value, args []string) ([]string, error) {
	var unnamed []string
	for i := 0; i < len(args); i++ {
		// collect non flag parameters
		if !strings.HasPrefix(args[i], "-") && args[i] != "-" {
			unnamed = append(unnamed, args[i])
			continue
		}
		// Locate field in struct of the flag name
		arg := strings.TrimLeft(args[i], "-")
		fld := reflection.FindFieldByName(arg, stc.Type().Elem(), reflection.FlagTag)
		if fld == nil {
			return nil, fmt.Errorf("--%s is an unknown flag", arg)
		}

		i++
		// create fld value from next argument or its default value
		var ival interface{}
		var err error
		if i < len(args) {
			ival, err = reflection.ValueFromString(args[i], fld.Type)
		}
		// If no valid value following flag, check if its an optional value flag.
		if ival == nil || err != nil {
			optVal := containsValue(reflection.TagOptionalValue, strings.Split(fld.Tag.Get(reflection.FlagTag), ","))
			if !optVal && fld.Type.Kind() != reflect.Bool {
				return nil, fmt.Errorf("missing value for flag -%s  %v", arg, err)
			}
			// Optional value with no value following, create default instance for field.
			ival, err = reflection.ValueFromString("", fld.Type)
			if err != nil {
				return nil, err
			}

			i-- // wind back arg as value not consumed
		}
		if err := reflection.SetFieldValue(stc.Elem().FieldByName(fld.Name), ival); err != nil {
			return nil, err
		}
	}
	return unnamed, nil
}

// structFromMethodFunc establishes the struct Type from the given method.
// Given mentod must be a Func which is a Method, (First parameter being the owning struct)
func structFromMethodFunc(i interface{}) (reflect.Type, error) {
	vt := reflect.TypeOf(i)
	if vt.Kind() != reflect.Func {
		return nil, fmt.Errorf("%s is not a method", vt.Name())
	}
	if vt.NumIn() < 1 {
		return nil, fmt.Errorf("%s is not a method", vt.Name())
	}
	st := vt.In(0)
	if st.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s is not a method", vt.Name())
	}
	return st, nil
}

// methodFromMethodFunc gets the method on the given struct instance, as named by the given func interface
func methodFromMethodFunc(i interface{}, st reflect.Type) (*reflect.Method, error) {
	fName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fns := strings.Split(fName, ".")
	fn := fns[len(fns)-1]
	m, ok := st.MethodByName(fn)
	if !ok {
		return nil, fmt.Errorf("method %s could not be found in struct %s", fn, st.Name())
	}
	return &m, nil
}

func containsValue(s string, values []string) bool {
	for _, v := range values {
		if v == s {
			return true
		}
	}
	return false
}

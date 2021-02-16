package flags

import (
	"errors"
	"fmt"
	"github.com/eurozulu/mainline/values"
	"reflect"
	"strings"
)

var ErrUnknownFlag = errors.New("unknown flag")

func NewFlags(ignoreUnknown bool) *Flags {
	return &Flags{flags: make(map[string]interface{}), IgnoreUnknown: ignoreUnknown}
}

type Flags struct {
	IgnoreUnknown bool
	flags         map[string]interface{}
}

func (fs Flags) Parameters() []string {
	v, ok := fs.flags[""]
	if !ok {
		return nil
	}
	return v.([]string)
}

func (fs Flags) String() []string {
	var ss []string
	for k, v := range fs.flags {
		if k == "" {
			ss = append(v.([]string), ss...)
			continue
		}
		ss = append(ss, fmt.Sprintf("-%s", k), fmt.Sprintf("\"%v\"", v))
	}
	return ss
}

func (fs Flags) Apply(args ...string) error {
	for i := 0; i < len(args); i++ {
		// collect non flag parameters in empty key
		if !strings.HasPrefix(args[i], "-") || args[i] == "-" {
			fs.flags[""] = append(fs.Parameters(), args[i])
			continue
		}

		arg := strings.TrimLeft(args[i], "-")
		v, ok := fs.flags[arg]
		if !ok {
			// unknown flag
			if fs.IgnoreUnknown {
				fs.flags[""] = append(fs.Parameters(), strings.Join([]string{"-", arg}, ""))
				continue
			}
			return fmt.Errorf("-%s is an %v", arg, ErrUnknownFlag)
		}

		var argVal string
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			i++
			argVal = args[i]
		}

		to := reflect.TypeOf(v)
		iVal, err := values.ValueFromString(argVal, to)
		if err != nil {
			// special case for bool.  If following arg not a bool "true" / "false", ignore it.
			if to.Kind() == reflect.Bool {
				iVal = true
				i--
			} else {
				return fmt.Errorf("could not read '%s' for flag -%s  %v", argVal, arg, err)
			}
		}
		values.SetValue(reflect.ValueOf(v), iVal)
	}
	return nil
}

// AddFlag adds one or more flag names which map to the given 'v' pointer to a variable
// at least one name must be given and it may optionall be followed by more names, all pointing to the same value.
// e.g. AddFlag(&Verbose, "verbose", "v")
// v must be a non nil pointer to a variable which will act as the receiver for the flag.
// If v is not a pointer, an error is thrown.  The pointer defines the data type of the flag,
// arguments following the flag, on the command line, will be parsed as that data type during apply.
func (fs *Flags) AddFlag(v interface{}, names ...string) error {
	if len(names) == 0 {
		return fmt.Errorf("flag name is missing")
	}

	if v == nil {
		return fmt.Errorf("flag value for '%s' is nil", strings.Join(names, " "))
	}
	val := reflect.ValueOf(v)
	if val.IsNil() {
		return fmt.Errorf("flag value for '%s' is nil", strings.Join(names, " "))
	}
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("flag value for '%s' is not a pointer", strings.Join(names, " "))
	}

	for _, n := range names {
		if _, ok := fs.flags[n]; ok {
			return fmt.Errorf("duplicate flag name.  '%s' already exists.", n)
		}
		fs.flags[n] = val.Interface()
	}
	return nil
}

// wildcardFlagMap attempts to find a Field in the given structure with a "Flags" tag option of wildcard "*".
// If a field is tagged as a wildcard flag, it must be defined as a map with string keys.
// Any flag not defined in the structure will be placed in the wildcard map.
// If no wildcard flag is set, Flags with no matching field throw the unknown flag error.
// Using a wildcard will prevent any error for unknown flag.
/*
func wildcardFlagMap(st reflect.Value) map[string]interface{} {
	// wildcard is optinal flag to collect undefined Flags
	wcfld := values.FindFieldByName(values.TagWildcard, st.Type().Elem(), values.FlagTag)
	if wcfld == nil {
		return nil
	}
	if wcfld.Type.Kind() != reflect.Map {
		log.Println("config error: wildcard flag field is not a map")
		return nil
	}
	fv := st.Elem().FieldByName(wcfld.Name)
	if fv.IsNil() {
		mp := reflect.MakeMapWithSize(wcfld.Type, 5)
		fv.Set(mp)
	}
	return fv.Interface().(map[string]interface{})
}
*/

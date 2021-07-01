package flags

import (
	"fmt"
	"github.com/eurozulu/commandgo/values"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

const TagName = "flag"
const TagHide = "-"

// Flags are arguments starting with a dash, usually followed by a value.  PArsed into global variables or struct fields.
// the key is the name of the flag, mapping to a value which must be a pointer
// Will panic if value is nil
type Flags map[string]interface{}

// GlobalFlags are defined prior to calling any command.  These map command line flags into global variables.
var GlobalFlags = Flags{}

// Apply applies the given command line flags to their respective, mapped positions.
// flags must begin with a '-' followed by a name mapped in this Flagmap (Case insensitive)
// All mapped flags have the argument following the flag parsed into the correct type for the mapped location,
// with the exception of booleans.  If the following arg can be parsed as a bool, it will be respected, otherwise
// the presents of a bool flag will set it to true.
// returns the given argument list with all of the 'used' flags and values removed. i.e. all non flags and unrecognised flags
func (f Flags) Apply(args ...string) ([]string, error) {
	var nonFlags []string

	for i := 0; i < len(args); i++ {
		// collect non flag parameters
		if !isFlagArg(args[i]) {
			nonFlags = append(nonFlags, args[i])
			continue
		}

		arg, ok := f.Contains(args[i])
		if !ok {
			// unknown flag, add to the non flags pile
			nonFlags = append(nonFlags, args[i])
			continue
		}
		v, ok := f[arg]

		var argVal string
		if i+1 < len(args) && !isFlagArg(args[i+1]) {
			i++
			argVal = args[i]
		}

		to := reflect.TypeOf(v)
		iVal, err := values.ValueFromString(argVal, to)
		if err != nil {
			// special case for bool.  If following arg not a bool "true" / "false", ignore it.
			if to.Kind() != reflect.Bool {
				return nil, fmt.Errorf("could not read '%s' for flag -%s  %v", argVal, arg, err)
			}
			iVal = true
			i--
		}
		values.SetValue(reflect.ValueOf(v), iVal)
	}
	return nonFlags, nil
}

// Checks if the given flag name is known.  Case insensitive
func (f Flags) Contains(arg string) (string, bool) {
	arg = strings.TrimLeft(arg, "-")
	for k := range f {
		if strings.EqualFold(k, arg) {
			return k, true
		}
	}
	return "", false
}

func isFlagArg(arg string) bool {
	return strings.HasPrefix(arg, "-") && arg != "-"
}

func fieldTagNames(fd *reflect.StructField) []string {
	tag, ok := fd.Tag.Lookup(TagName)
	if !ok { // no tag
		return nil
	}
	var names []string
	for _, tn := range strings.Split(tag, ",") {
		names = append(names, strings.TrimSpace(tn))
	}
	return names
}

// NewFlagFromStruct creates a new Flags for the fields in the given struct type.
func NewFlagsFromStruct(str reflect.Value) Flags {
	if str.Kind() != reflect.Ptr && str.Elem().Kind() != reflect.Struct {
		log.Fatalln(fmt.Errorf("struct flags failed as given item not a ptr to structure"))
	}

	fs := Flags{}
	v := str.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if _, ok := fld.Tag.Lookup(TagHide); ok {
			continue
		}
		names := fieldTagNames(&fld)
		if len(names) == 0 {
			names = []string{strings.ToLower(fld.Name)}
		}
		// Convert field into a ptr of its type
		vfld := v.Field(i)
		pFld := reflect.NewAt(vfld.Type(), unsafe.Pointer(vfld.UnsafeAddr()))
		for _, name := range names {
			fs[name] = pFld.Interface()
		}
	}
	return fs
}

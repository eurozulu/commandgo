package argdecode

import (
	"fmt"
	"github.com/eurozulu/flinger/logger"
	"reflect"
	"strconv"
	"strings"
)

const tagFlag = "flag"
const tagCommand = "command"

type Decoder struct {
	Args []string
}

func (d Decoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer passed to Unmarshal")
	}
	if val.Kind() == reflect.Interface && !val.IsNil() {
		e := val.Elem()
		if e.Kind() == reflect.Ptr && !e.IsNil() {
			val = e
		}
	}
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	return d.unmarshal(val)
}

func (d *Decoder) unmarshal(val reflect.Value) error {
	var unnamed []string
	for i := 0; i < len(d.Args); i++ {
		if !strings.HasPrefix(d.Args[i], "-") && d.Args[i] != "-" {
			unnamed = append(unnamed, d.Args[i])
			continue
		}

		// Locate field in struct of the flag name
		arg := strings.TrimLeft(d.Args[i], "-")
		fld := findFieldByName(arg, val, tagFlag)
		if fld == nil {
			return fmt.Errorf("--%s is an unknown flag", arg)
		}

		// Get next arg as value and convert to correct type
		var s string
		if (i + 1) < len(d.Args) {
			i++
			s = d.Args[i]
			// Check if a bool type and not a bool value, assume next arg not related.
			if fld.Type.Kind() == reflect.Bool {
				_, err := strconv.ParseBool(s)
				if err != nil {
					i--
					s = ""
				}
			}
		}
		iv, err := ValueFromString(s, fld.Type)
		if err != nil {
			return err
		}

		if err := setFieldValue(val.FieldByName(fld.Name), iv); err != nil {
			return err
		}
	}

	if len(unnamed) == 0 {
		return nil
	}

	// Check if any field is taged as 'command' matching first unnamed
	cmdft := findFieldByName(unnamed[0], val, tagCommand)
	if cmdft != nil {
		logger.Debug("command %d mapped to function", unnamed[0])
		if err := callCommandFunc(cmdft, val, unnamed); err != nil {
			return err
		}
	}

	return nil
}

// callCommandFunc calls the function which is the value of the cmdfld, using the given args.
// Args are translated into the correct types for the call.
func callCommandFunc(cmdfld *reflect.StructField, val reflect.Value, args []string) error {
	cmdfv := val.FieldByName(cmdfld.Name)
	if cmdfv.IsNil() {
		return fmt.Errorf("ignoring command %s as no function has been set for that command", args[0])
	}

	if cmdfv.Type().Kind() != reflect.Func {
		return fmt.Errorf("failed to start %s as Field %s is not a function", args[0], cmdfld.Name)
	}
	isMth := IsMethod(cmdfv)
	sig := NewSignature(cmdfv.Type(), isMth)
	if len(sig.ParamTypes) != (len(args) - 1) {
		return fmt.Errorf("%s requires %d arguments.  Found %d", args[0], len(sig.ParamTypes), len(args)-1)
	}

	iVals, err := ValuesFromString(args, sig.ParamTypes)
	if err != nil {
		return fmt.Errorf("failed to read arguments for %s  %v", args[0], err)
	}
	vals := make([]reflect.Value, len(iVals))
	for i, iv := range iVals {
		vals[i] = reflect.ValueOf(iv)
	}
	val.FieldByName(cmdfld.Name).Elem().Call(vals)
	return nil
}

// findFieldByName scans each field in the given struct for either its fieldname or one of its 'flag' tag names, for the given name.
func findFieldByName(name string, str reflect.Value, tagName string) *reflect.StructField {
	for i := 0; i < str.NumField(); i++ {
		fld := str.Type().Field(i)
		names := fieldNames(fld, tagName)
		for _, n := range names {
			if strings.EqualFold(name, n) {
				return &fld
			}
		}
	}
	return nil
}

// Gets the names of the given field. Includes the field name and any comma separated names found in the given tag.
func fieldNames(fd reflect.StructField, tagName string) []string {
	var names = []string{fd.Name}
	tag, ok := fd.Tag.Lookup(tagName)
	if !ok { // no tag, just the field name
		return names
	}
	if tag == "-" { // ignore those taged with dash
		return nil
	}
	for _, tn := range strings.Split(tag, ",") {
		names = append(names, strings.TrimSpace(tn))
	}
	return names
}

// Sets the given field value, the given value.
// Assigns the value or a pointer to it, depending on the field type
func setFieldValue(fld reflect.Value, val interface{}) error {
	var vp reflect.Value
	v := reflect.ValueOf(val)
	if v.Type().Kind() == reflect.Ptr {
		vp = v
		v = v.Elem()
	}
	// Assign field value.  Check if receiver is expecting a pointer or not.
	if fld.Type().Kind() == reflect.Ptr {
		fld.Set(vp)
	} else {
		fld.Set(v)
	}
	return nil
}

func NewDecoder(args []string) *Decoder {
	return &Decoder{Args: args}
}

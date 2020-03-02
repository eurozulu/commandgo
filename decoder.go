package mainline

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Decoder decodes command line arguments into objects
type Decoder struct {
	Args []string
}

// Decode the arguments into the given object.
// Arguments are parsed according to the structure of the object and its types.
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

	// process the unnamed 'parameters', assign to flag marked'*'
	argft := findFieldByName("*", val, tagFlag)
	if len(unnamed) == 0 {
		if argft == nil {
			return nil
		}
		return fmt.Errorf("expected arguments for %s", argft.Name)
	}
	if argft == nil {
		return fmt.Errorf("unexpected arguments %v", unnamed)
	}

	args := strings.Join(unnamed, SliceDelimiter)
	v, err := ValueFromString(args, argft.Type)
	if err != nil {
		return fmt.Errorf("failed to parse arguments %v", err)
	}
	if err := setFieldValue(val.FieldByName(argft.Name), v); err != nil {
		return fmt.Errorf("failed to assign arguments %v to field %s  %v", v, argft.Name, err)
	}
	return nil
}

func NewDecoder(args []string) *Decoder {
	return &Decoder{Args: args}
}

package mainline

import (
	"reflect"
	"strings"
)

const tagFlag = "flag"

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

package flags

import (
	"fmt"
	"reflect"
	"strings"
)

const TagName = "flag"
const TagHide = "-"
const TagOptionalValue = "optionalvalue"
const TagWildcard = "*"

func NewStructFlags(str reflect.Type) (*Flags, error) {
	if str.Kind() == reflect.Ptr {
		return NewStructFlags(str.Elem())
	}
	if str.Kind() != reflect.Struct {
		return nil, fmt.Errorf("struct flags failed as given item not a structure")
	}

	fs := NewFlags(false)
	for _, fld := range structFields(str) {
		// Use tag names for flag names unles no tags available, then use field name in lower case
		names := fieldTagNames(fld)
		if len(names) == 0 {
			names = []string{strings.ToLower(fld.Name)}
		}
		if err := fs.AddFlag(fld, names...); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

func structFields(t reflect.Type) []*reflect.StructField {
	var sfs []*reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if _, ok := fld.Tag.Lookup(TagHide); ok {
			continue
		}

		// If struct field, recurse down through this
		if fld.Type.Kind() == reflect.Struct {
			sfs = append(sfs, structFields(fld.Type)...)
			continue
		}
		sfs = append(sfs, &fld)
	}
	return sfs
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

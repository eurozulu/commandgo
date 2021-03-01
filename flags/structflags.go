// Copyright 2020 Rob Gilham
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package flags

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

const TagName = "flag"
const TagHide = "-"

//const TagOptionalValue = "optionalvalue"
//const TagWildcard = "*"

// NewStructFlags creates a new Flags parser for the fields in the given struct type.
func NewStructFlags(str reflect.Value) (*Flags, error) {
	if str.Kind() != reflect.Ptr && str.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("struct flags failed as given item not a ptr to structure")
	}

	fs := NewFlags(false)

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
		if err := fs.AddFlag(pFld.Interface(), names...); err != nil {
			return nil, err
		}
	}
	return fs, nil
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

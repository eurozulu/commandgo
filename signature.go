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

package mainline

import (
	"bytes"
	"fmt"
	"reflect"
)

// Signature represents the signature of a method, both its parameters and its return types.
type Signature struct {
	ParamTypes  []reflect.Type
	ReturnTypes []reflect.Type
}

func (s Signature) String() string {
	ps := s.listTypes(s.ParamTypes)
	if len(s.ReturnTypes) > 0 {
		ps = fmt.Sprintf("%s  Returns %s", ps, s.listTypes(s.ReturnTypes))
	}
	return ps
}

func (s Signature) listTypes(t []reflect.Type) string {
	bf := bytes.NewBuffer(nil)
	for i, p := range t {
		if i > 0 {
			bf.WriteString(", ")
		}
		bf.WriteString(p.String())
	}
	return fmt.Sprintf("[%s]", bf.String())
}

// NewSignature creates a new signature from the given method value
func NewSignature(t reflect.Type, isMethod bool) Signature {
	start := 0
	if isMethod {
		start++
	}

	params := make([]reflect.Type, t.NumIn()-start)
	x := 0
	for i := start; i < t.NumIn(); i++ {
		params[x] = t.In(i)
		x++
	}

	returns := make([]reflect.Type, t.NumOut())
	for i := 0; i < t.NumOut(); i++ {
		returns[i] = t.Out(i)
	}
	return Signature{
		ParamTypes:  params,
		ReturnTypes: returns,
	}
}

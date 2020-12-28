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

package reflection

import (
	"bytes"
	"fmt"
	"reflect"
)

// Signature represents the signature of a method or func, both its parameters and its return types.
type Signature struct {
	ParamTypes  []reflect.Type
	ReturnTypes []reflect.Type
	IsVariadic  bool
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

// NewSignatureOf creates a Signature of the given func or Method
func NewSignatureOf(fun interface{}) (*Signature, error) {
	var t reflect.Type
	m, ok := fun.(reflect.Method)
	if ok {
		t = m.Type

	} else { // not a method, see if its a func
		t = reflect.TypeOf(fun)
		if t.Kind() != reflect.Func {
			return nil, fmt.Errorf("fun %s is not a method or func", t.Name())
		}
	}

	var index int
	if ok { // if method, move pased first param
		index++
	}
	in := t.NumIn()
	var params []reflect.Type
	for ; index < in; index++ {
		params = append(params, t.In(index))
	}
	out := t.NumOut()
	returns := make([]reflect.Type, out)
	for i := 0; i < out; i++ {
		returns[i] = t.Out(i)
	}
	return &Signature{
		ParamTypes:  params,
		ReturnTypes: returns,
		IsVariadic:  t.IsVariadic(),
	}, nil
}

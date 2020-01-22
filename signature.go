package mainline

import (
	"bytes"
	"fmt"
	"reflect"
)

// Signature represents the signature of a method, both its parameters and its return types.
type Signature struct {
	Method      reflect.Method
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
func NewSignature(mv reflect.Method) *Signature {
	t := mv.Type
	params := make([]reflect.Type, t.NumIn()-1)
	for i := 1; i < t.NumIn(); i++ {
		params[i-1] = t.In(i)
	}
	returns := make([]reflect.Type, t.NumOut())
	for i := 0; i < t.NumOut(); i++ {
		returns[i] = t.Out(i)
	}
	return &Signature{
		Method:      mv,
		ParamTypes:  params,
		ReturnTypes: returns,
	}
}

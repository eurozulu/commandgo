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
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func ValueToString(v interface{}) string {
	vt := reflect.TypeOf(v)

	switch vt.Kind() {
	case reflect.Ptr:
		return ValueToString(reflect.ValueOf(v).Elem().Interface())

	case reflect.Struct:
		return structureToString(v)

	case reflect.Slice:
		return valueToJsonString(v)

	case reflect.Map:
		return valueToJsonString(v)

	case reflect.Float64, reflect.Float32:
		return strconv.FormatFloat(v.(float64), 64, 2, 64)

	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return strconv.FormatInt(v.(int64), 10)

	case reflect.Bool:
		return strconv.FormatBool(v.(bool))

	case reflect.String:
		return v.(string)

	default:
		return ""
	}
}

func valueToJsonString(v interface{}) string {
	by, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(by)
}

func structureToString(v interface{}) string {
	tm, ok := v.(encoding.TextMarshaler)
	if ok {
		by, err := tm.MarshalText()
		if err != nil {
			return ""
		}
		return string(by)
	}
	jm, ok := v.(json.Marshaler)
	if ok {
		by, err := jm.MarshalJSON()
		if err != nil {
			return ""
		}
		return string(by)
	}

	return fmt.Sprintf("%vv", v)
}

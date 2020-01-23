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

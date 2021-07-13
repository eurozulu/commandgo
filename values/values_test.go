package values_test

import (
	"commandgo/values"
	"log"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type testIntType int
type testStringType string

var testvarBool bool
var testvarString string
var testvarUrl *url.URL

func TestIsKind(t *testing.T) {
	if !values.IsKind(testvarBool, reflect.Bool) {
		t.Fatalf("unexpected kind checking bool")
	}
	if values.IsKind(testvarString, reflect.Bool) {
		t.Fatalf("unexpected kind checking bool")
	}
	var i interface{}
	if values.IsKind(i, reflect.Bool) {
		t.Fatalf("unexpected kind checking bool")
	}

	if !values.IsKind(testvarUrl, reflect.Struct) {
		t.Fatalf("unexpected kind checking nil url")
	}
}

func TestValueFromString_int(t *testing.T) {
	v, err := values.ValueFromString("555", reflect.TypeOf(0))
	if err != nil {
		t.Fatalf("%v", err)
	}
	vv, ok := v.(int)
	if !ok {
		t.Fatalf("unexpected type found, expected int, found %v", reflect.TypeOf(v))
	}
	if vv != 555 {
		t.Fatalf("unexpected value found, expected 555, found %v", reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("0", reflect.TypeOf(0))
	if err != nil {
		t.Fatalf("%v", err)
	}
	vv, ok = v.(int)
	if !ok {
		t.Fatalf("unexpected type found, expected int, found %v", reflect.TypeOf(v))
	}
	if vv != 0 {
		t.Fatalf("unexpected value found, expected 0, found %v", reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("-123", reflect.TypeOf(0))
	if err != nil {
		t.Fatalf("%v", err)
	}
	vv, ok = v.(int)
	if !ok {
		t.Fatalf("unexpected type found, expected int, found %v", reflect.TypeOf(v))
	}
	if vv != -123 {
		t.Fatalf("unexpected value found, expected -123, found %v", reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("9999999999", reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("%v", err)
	}
	vvv, ok := v.(int64)
	if !ok {
		t.Fatalf("unexpected type found, expected int64, found %v", reflect.TypeOf(v))
	}
	if vvv != 9999999999 {
		t.Fatalf("unexpected value found, expected 0, found %v", reflect.ValueOf(v))
	}

	// Attempt an overflow value
	v, err = values.ValueFromString("9999999999", reflect.TypeOf(int16(0)))
	if err == nil || !strings.HasSuffix(err.Error(), "argument 9999999999 could not be parsed as a int16") {
		t.Fatalf("expected error parsing overflow value to int16 %v", err)
	}
	v, err = values.ValueFromString("255", reflect.TypeOf(uint8(0)))
	if err != nil {
		t.Fatalf("unexpected error parsing uint8%v", err)
	}
	vvvv, ok := v.(uint8)
	if !ok {
		t.Fatalf("unexpected type found, expected uint8, found %v", reflect.TypeOf(v))
	}
	if vvvv != 255 {
		t.Fatalf("unexpected value found, expected 0, found %v", reflect.ValueOf(v))
	}

	var tt testIntType = 0
	v, err = values.ValueFromString("255", reflect.TypeOf(tt))
	if err != nil {
		t.Fatalf("unexpected error parsing custom int type %v", err)
	}
	ttt, ok := v.(testIntType)
	if !ok {
		t.Fatalf("unexpected return value type parsing custom int type. Expected %v, found %v", reflect.TypeOf(tt), reflect.TypeOf(v))
	}
	if ttt != 255 {
		t.Fatalf("unexpected value parsing custom int type. Expected %d, found %d", 255, tt)
	}
}

func TestValueFromString_float(t *testing.T) {
	v, err := values.ValueFromString("1.23", reflect.TypeOf(float64(0)))
	if err != nil {
		t.Fatalf("unexpected error parsing float value %v", err)
	}
	vv, ok := v.(float64)
	if !ok {
		t.Fatalf("unexpected type found, expected float64, found %v", reflect.TypeOf(v))
	}
	if vv != 1.23 {
		t.Fatalf("unexpected value found, expected 1.23, found %v", reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("1.23", reflect.TypeOf(float32(0)))
	if err != nil {
		t.Fatalf("unexpected error parsing float value %v", err)
	}
	vvv, ok := v.(float32)
	if !ok {
		t.Fatalf("unexpected type found, expected float32, found %v", reflect.TypeOf(v))
	}
	if vvv != 1.23 {
		t.Fatalf("unexpected value found, expected 1.23, found %v", reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("0", reflect.TypeOf(float32(0)))
	if err != nil {
		t.Fatalf("unexpected error parsing float value %v", err)
	}
	vvv, ok = v.(float32)
	if !ok {
		t.Fatalf("unexpected type found, expected float32, found %v", reflect.TypeOf(v))
	}
	if vvv != 0 {
		t.Fatalf("unexpected value found, expected 0, found %v", reflect.ValueOf(v))
	}
}

func TestValueFromString_bool(t *testing.T) {
	v, err := values.ValueFromString("true", reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("%v", err)
	}
	vv, ok := v.(bool)
	if !ok {
		t.Fatalf("unexpected type found, expected bool, found %v", reflect.TypeOf(v))
	}
	if vv != true {
		t.Fatalf("unexpected value found, expected bool true, found %v", reflect.ValueOf(v))
	}

	// test with no param
	v, err = values.ValueFromString("", reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("%v", err)
	}
	vv, ok = v.(bool)
	if !ok {
		t.Fatalf("unexpected type found, expected bool, found %v", reflect.TypeOf(v))
	}
	if vv != true {
		t.Fatalf("unexpected value found, expected bool true, found %v", reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("false", reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("%v", err)
	}
	vv, ok = v.(bool)
	if !ok {
		t.Fatalf("unexpected type found, expected bool, found %v", reflect.TypeOf(v))
	}
	if vv != false {
		t.Fatalf("unexpected value found, expected bool false, found %v", reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("blabla", reflect.TypeOf(true))
	if err == nil || !strings.HasSuffix(err.Error(), "blabla could not be read as a bool") {
		t.Fatalf("expected error with bad bool value")
	}
}

func TestValueFromString_string(t *testing.T) {
	testString := "Hello world"
	v, err := values.ValueFromString(testString, reflect.TypeOf(""))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	vv, ok := v.(string)
	if !ok {
		t.Fatalf("unexpected type found, expected string, found %v", reflect.TypeOf(v))
	}
	if vv != testString {
		t.Fatalf("unexpected value found, expected %s, found %v", testString, reflect.ValueOf(v))
	}

	v, err = values.ValueFromString("", reflect.TypeOf(""))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	vv, ok = v.(string)
	if !ok {
		t.Fatalf("unexpected type found, expected string, found %v", reflect.TypeOf(v))
	}
	if vv != "" {
		t.Fatalf("unexpected value found, expected %s, found %v", "", reflect.ValueOf(v))
	}

	var ct testStringType
	v, err = values.ValueFromString(testString, reflect.TypeOf(ct))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	if reflect.TypeOf(v).Name() != reflect.TypeOf(ct).Name() {
		t.Fatalf("unexpected type found, expected string, found %v", reflect.TypeOf(v))
	}
	if v.(testStringType) != testStringType(testString) {
		t.Fatalf("unexpected value found, expected %s, found %v", testString, reflect.ValueOf(v))
	}
}

func TestValueFromString_map(t *testing.T) {
	testString := "{\"one\": 1, \"two\": 2, \"three\": 3}"

	v, err := values.ValueFromString(testString, reflect.TypeOf(map[string]int{}))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	vv, ok := v.(map[string]int)
	if !ok {
		t.Fatalf("unexpected type found, expected map[string] int, found %v", reflect.TypeOf(v))
	}
	if len(vv) != 3 {
		t.Fatalf("unexpected map length found, expected %d, found %v", 3, len(vv))
	}
	i, ok := vv["one"]
	if !ok {
		log.Fatalln("expected value 'one' not found in returned map")
	}
	if i != 1 {
		log.Fatalf("unexpected value 'one' found in returned map. expected %d, found %d", 1, i)
	}
	i, ok = vv["two"]
	if !ok {
		log.Fatalln("expected value 'two' not found in returned map")
	}
	if i != 2 {
		log.Fatalf("unexpected value 'two' found in returned map. expected %d, found %d", 2, i)
	}
	i, ok = vv["three"]
	if !ok {
		log.Fatalln("expected value 'three' not found in returned map")
	}
	if i != 3 {
		log.Fatalf("unexpected value 'three' found in returned map. expected %d, found %d", 3, i)
	}
}

func TestValueFromString_slice(t *testing.T) {
	testString := "one, two, three"
	v, err := values.ValueFromString(testString, reflect.TypeOf([]string{}))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	vv, ok := v.([]string)
	if !ok {
		t.Fatalf("unexpected type found, expected []string, found %v", reflect.TypeOf(v))
	}
	if len(vv) != 3 {
		t.Fatalf("unexpected slice length found, expected %d, found %v", 3, len(vv))
	}
	if vv[0] != "one" || vv[1] != "two" || vv[2] != "three" {
		t.Fatalf("unexpected values returned in slice")
	}

	testString = "1, 2, 3"
	v, err = values.ValueFromString(testString, reflect.TypeOf([]int64{}))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	vvv, ok := v.([]int64)
	if !ok {
		t.Fatalf("unexpected type found, expected []int64, found %v", reflect.TypeOf(v))
	}
	if len(vvv) != 3 {
		t.Fatalf("unexpected slice length found, expected %d, found %v", 3, len(vv))
	}
	for i, iv := range vvv {
		if int64(i)+1 != iv {
			t.Fatalf("unexpected values returned in int64 slice")
		}
	}

	testString = "1, 2, hello"
	v, err = values.ValueFromString(testString, reflect.TypeOf([]int{}))
	if err == nil || !strings.HasSuffix(err.Error(), "hello could not be read as a int") {
		t.Fatalf("expected error parsing string value %v", err)
	}
}

type testjsonstruct struct {
	One   string       `json:"one,omitempty"`
	Two   float64      `json:"two,omitempty"`
	Three int          `json:"three,omitempty"`
	Name  *testsubtype `json:"name,omitempty"`
}
type testsubtype struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
}

func TestValueFromString_struct(t *testing.T) {
	testString := "{\"one\": \"hello\", \"two\": 1.23, \"three\": 555, \"name\": { \"first\": \"john\", \"last\": \"doe\"} }"
	v, err := values.ValueFromString(testString, reflect.TypeOf(&testjsonstruct{}))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	vv, ok := v.(*testjsonstruct)
	if !ok {
		t.Fatalf("unexpected type found, expected *testjsonstruct, found %v", reflect.TypeOf(v))
	}
	if vv.One != "hello" || vv.Two != 1.23 || vv.Three != 555 {
		t.Fatalf("unexpected values found in returned test object")
	}
	if vv.Name == nil {
		t.Fatalf("unexpected nil values found in returned test object")
	}
	if vv.Name.First != "john" || vv.Name.Last != "doe" {
		t.Fatalf("unexpected values found in returned test object")
	}
}

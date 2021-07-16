package values_test

import (
	"fmt"
	"github.com/eurozulu/commandgo/values"
	"net/url"
	"reflect"
	"testing"
)

func TestCustomValueFromString(t *testing.T) {
	testString := "http://www.google.com"
	v, err := values.ValueFromString(testString, reflect.TypeOf(&url.URL{}))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	vv, ok := v.(*url.URL)
	if !ok {
		t.Fatalf("unexpected type found, expected *url.URL, found %v", reflect.TypeOf(v))
	}
	if vv.String() != testString {
		t.Fatalf("unexpected return url testing custom type. expected %s, found %s", testString, vv.String())
	}
}

type UserID string

func TestCustomValueFromString_type(t *testing.T) {
	testdb := map[string]UserID{
		"a": "12345",
		"b": "67890",
		"c": "abcde",
	}

	values.NewCustomType(reflect.TypeOf(UserID("")), func(s string, t reflect.Type) (interface{}, error) {
		id, ok := testdb[s]
		if !ok {
			return nil, fmt.Errorf("user %s unknown", s)
		}
		return id, nil
	})

	v, err := values.ValueFromString("a", reflect.TypeOf(UserID("")))
	if err != nil {
		t.Fatalf("unexpected error parsing string value %v", err)
	}
	if reflect.TypeOf(v).Name() != reflect.TypeOf(UserID("")).Name() {
		t.Fatalf("unexpected type returned with custom type check. expected UserID, found %v", reflect.TypeOf(v))
	}
	vv := v.(UserID)
	if vv != UserID("12345") {
		t.Fatalf("unexpected value returned with custom type check. expected %s, found %s", "12345", vv)
	}
}

func testUserId(id UserID) error {
	return nil
}

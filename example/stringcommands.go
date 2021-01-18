package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
)

type StringCommands struct {
	// WrapText -wrap or -w wraps the given string before and after the reversed result.
	WrapText string `flag:"w,wrap"`

	// MapTest tagged as the wildcard flag with '*'.  Any flags specified not matching to fields are placed in this map.
	MapTest    map[string]interface{} `flag:"map,*"`
	URLTest    *url.URL               `flag:"url,optionalvalue"`
	StrPtrTest *string                `flag:"str,optionalvalue"`
	IntPtrTest *int                   `flag:"int,optionalvalue"`
}

func (sc StringCommands) Test() {

	if sc.IntPtrTest == nil {
		fmt.Println("int nil")
	} else {
		fmt.Printf("int: %d\n", *sc.IntPtrTest)
	}

	if sc.URLTest == nil {
		fmt.Println("url nil")
	} else {
		fmt.Printf("url: %v\n", sc.URLTest.String())
	}

	if sc.StrPtrTest == nil {
		fmt.Println("StrPtrTest nil")
	} else {
		fmt.Printf("StrPtrTest: '%s'\n", *sc.StrPtrTest)
	}

	if sc.MapTest == nil {
		fmt.Println("Map not set")
		return
	}
	if len(sc.MapTest) == 0 {
		sc.MapTest["map"] = "is empty!"
	}
	for k, v := range sc.MapTest {
		fmt.Printf("%s = %v", k, v)
	}

	fmt.Println("\ndone")
}

// Reverse the given argument
func (sc StringCommands) Reverse(s string) {
	b := bytes.NewBuffer(nil)
	if sc.WrapText != "" {
		b.WriteString(sc.WrapText)
	}
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteRune(rune(s[i]))
	}
	if sc.WrapText != "" {
		b.WriteString(sc.WrapText)
	}
	fmt.Println(b.String())
}

func (sc StringCommands) Base64Encode(s string) {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(s)))
}

func (cr StringCommands) Square(i1, i2 int) {
	fmt.Println(i1 ^ i2)
}

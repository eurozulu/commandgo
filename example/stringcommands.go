package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

type StringCommands struct {
	// WrapText -wrap or -w wraps the given string before and after the reversed result.
	WrapText string `flag:"w,wrap"`
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

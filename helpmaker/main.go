package main

import (
	"github.com/eurozulu/mainline/helpmaker/helper"
)

func main() {
	by, err := helper.GenHelp()
	if err != nil {
		panic(err)
	}

	if err := helper.GenCode(by); err != nil {
		panic(err)
	}

}

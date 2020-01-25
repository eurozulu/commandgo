package main

import (
	"github.com/eurozulu/mainline/helpmaker/helper"
)

func main() {
	if err := helper.GenHelp(); err != nil {
		panic(err)
	}

	if err := helper.GenCode(); err != nil {
		panic(err)
	}

	if err := helper.CleanJson(); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"github.com/eurozulu/mainline"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	listCmd := &FileCommands{root: wd}

	cmds := mainline.Commands{
		"reverse":                    &StringCommands{},
		"square":                     &StringCommands{},
		"-Base64Encode, encode, enc": &StringCommands{},
		"listFiles, lf":              listCmd,
		"listdirectory, ld":          listCmd,
	}

	out, err := cmds.Run(os.Args...)
	if err != nil {
		fmt.Println(err)
	}
	for _, o := range out {
		fmt.Printf("%v\n", o)
	}
}

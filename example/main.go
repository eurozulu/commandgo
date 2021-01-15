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
		"test":                       &StringCommands{},
		"reverse":                    &StringCommands{},
		"square":                     &StringCommands{},
		"-Base64Encode, encode, enc": &StringCommands{},
		"listFiles, lf":              listCmd,
		"listdirectory, ld":          listCmd,
	}

	if err := cmds.Run(os.Args...); err != nil {
		fmt.Println(err)
	}
}

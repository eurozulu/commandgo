package main

import (
	"fmt"
	"github.com/eurozulu/mainline"
	"os"
)

func main() {
	cmds := mainline.Commands{
		"test":          StringCommands.Test,
		"reverse":       StringCommands.Reverse,
		"square":        StringCommands.Square,
		"Base64Encode":  StringCommands.Base64Encode,
		"encode":        StringCommands.Base64Encode,
		"listFiles":     FileCommands.ListFiles,
		"lf":            FileCommands.ListFiles,
		"ls":            FileCommands.ListFiles,
		"listdirectory": FileCommands.ListDirectory,
		"ld":            FileCommands.ListDirectory,
		"help":          mainline.HelpCommand.Help,
	}

	if err := cmds.Run(os.Args...); err != nil {
		fmt.Println(err)
	}
}

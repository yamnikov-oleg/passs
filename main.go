package main

import (
	"fmt"
	"os"
)

func PrintNoCmd() {
	fmt.Printf("Type '%s help' to get instructions\n", os.Args[0])
}

func main() {
	if len(os.Args) == 1 {
		PrintNoCmd()
		os.Exit(1)
	}

	var selCmd *Command
OuterLoop:
	for _, cmd := range cmds {
		for _, lbl := range cmd.Labels {
			if lbl == os.Args[1] {
				selCmd = &cmd
				break OuterLoop
			}
		}
	}

	if selCmd == nil {
		PrintNoCmd()
		os.Exit(1)
	}

	selCmd.Func(os.Args[2:])
}

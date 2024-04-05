package main

import (
	"fmt"
	"os"
)

func main() {
	var args []string
	var flags []string

	for _, arg := range os.Args[1:] {
		if len(arg) > 1 && arg[0] == '-' {
			flags = append(flags, arg[1:])
			continue
		}
		args = append(args, arg)
	}

	fmt.Println("Arguments:")
	for _, arg := range args {
		fmt.Println(arg)
	}

	fmt.Println("Flags:")
	for _, flag := range flags {
		fmt.Println(flag)
	}
}

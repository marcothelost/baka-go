package main

import (
	"fmt"
	"os"
)

const (
	TIMETABLE_COMMAND = "timetable"
	AVERAGES_COMMAND = "averages"
)

func timetable() {
	fmt.Println("Timetable")
}

func averages() {
	fmt.Println("Averages")
}

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

	if len(args) == 0 {
		fmt.Println("No arguments!")
		os.Exit(1)
	}

	switch (args[0]) {
	case TIMETABLE_COMMAND:
		timetable()
	case AVERAGES_COMMAND:
		averages()
	default:
		fmt.Println("Command not found")
		os.Exit(1)
	}
}

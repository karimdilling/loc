package options

import (
	"flag"
	"fmt"
)

var Verbose bool
var Help bool

func init() {
	flag.BoolVar(&Help, "h", false, "Prints help output")
	flag.BoolVar(&Help, "help", false, "Prints help output")

	flag.BoolVar(&Verbose, "v", false, "Prints the line numbers for every single source file")
	flag.BoolVar(&Verbose, "verbose", false, "Prints the line numbers for every single source file")
}

func PrintProgramDescription() {
	const programDescription = "Usage: loc [OPTION]... [FILE OR FOLDER]...\n" +
		"Use '.' as [FILE OR FOLDER] for the current directory\n" +
		"\n" +
		"This program prints the total and the effective line numbers of the specified files or whole directories " +
		"and does so only for valid source files (also ignores dotfiles/-folders as well as well known build directories)\n" +
		"The amount of effective lines is the number of total lines without the blank lines and comments specific to the according programming language.\n" +
		"\n" +
		"The following [OPTION]s can be set:"

	fmt.Println(programDescription)
}

package options

import "flag"

var Verbose bool
var Help bool

func init() {
	flag.BoolVar(&Verbose, "v", false, "Prints the line numbers for every single source file")
	flag.BoolVar(&Help, "h", false, "This program prints out the total line numbers of files or whole directories and only includes source files. The effective line count is a line count that dismisses blank lines and comments")
}

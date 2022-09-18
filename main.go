package main

import (
	"flag"
	"fmt"
	"loc/options"
	"loc/parse"
	"os"
)

func main() {
	flag.Parse()
	if options.Help {
		flag.PrintDefaults()
		return
	}

	filenames, err := parse.GetFilenames()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesAndLines := make(map[string][2]int)
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		parse.CountLines(file, filesAndLines)
		defer file.Close()
	}

	printLineNumbers(filesAndLines)
}

func printLineNumbers(filesAndLines map[string][2]int) {
	totalCount := 0
	effectiveTotalCount := 0
	for filename, count := range filesAndLines {
		if options.Verbose {
			fmt.Printf("%v: total %v, effective %v\n", filename, count[0], count[1])
		}
		totalCount += count[0]
		effectiveTotalCount += count[1]
	}

	fmt.Println("The total line count is", totalCount)
	fmt.Println("The effective line count is", effectiveTotalCount)
}

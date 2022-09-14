package parse

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Counts overall lines and the effective source code lines with comments and blank lines stripped out
func CountLines(file *os.File, fileMap map[string][2]int) {
	scanner := bufio.NewScanner(file)
	lineCount := 0
	effectiveCount := 0
	skipLine := false
	skipMultipleLines := false

	cLikeComments := strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), ".c") ||
		strings.HasSuffix(file.Name(), ".cpp") || strings.HasSuffix(file.Name(), ".cc") ||
		strings.HasSuffix(file.Name(), ".h") || strings.HasSuffix(file.Name(), ".hpp") ||
		strings.HasSuffix(file.Name(), ".java") || strings.HasSuffix(file.Name(), ".js") ||
		strings.HasSuffix(file.Name(), ".jsx") || strings.HasSuffix(file.Name(), ".ts") ||
		strings.HasSuffix(file.Name(), ".rs")
	phpComments := strings.HasSuffix(file.Name(), ".php")
	pythonComments := strings.HasSuffix(file.Name(), ".py")
	rubyComments := strings.HasSuffix(file.Name(), ".rb")
	rComments := strings.HasSuffix(file.Name(), ".R")

	for scanner.Scan() {
		lineCount++

		text := scanner.Text()
		text = strings.TrimSpace(text)

		if cLikeComments {
			parseCFile(&text, &effectiveCount, &skipLine, &skipMultipleLines)
		} else if phpComments {
			parsePHPFile(&text, &effectiveCount, &skipLine, &skipMultipleLines)
		} else if pythonComments {
			parsePythonFile(&text, &effectiveCount, &skipLine, &skipMultipleLines)
		} else if rubyComments {
			parseRubyFile(&text, &effectiveCount, &skipLine, &skipMultipleLines)
		} else if rComments {
			parseRFile(&text, &effectiveCount, &skipLine)
		}

		if !skipLine {
			effectiveCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("%s in file: %v\n", err, file.Name())
	}

	lineCountValues := [2]int{lineCount, effectiveCount}
	fileMap[file.Name()] = lineCountValues
}

func parseRFile(text *string, effectiveCount *int, skipLine *bool) {
	if strings.HasPrefix(*text, "#") {
		*skipLine = true
	} else if *text == "" {
		*skipLine = true
	} else {
		*skipLine = false
	}
}

func parseRubyFile(text *string, effectiveCount *int, skipLine *bool, skipMultipleLines *bool) {
	if *skipMultipleLines {
		if strings.HasPrefix(*text, "=end") {
			*skipMultipleLines = false
			if strings.HasSuffix(*text, "=end") { // check if code is following after ending the comment
				*skipLine = true
			} else {
				*skipLine = false
			}
		} else {
			*skipLine = true
			*skipMultipleLines = true
		}
		return
	}

	if strings.HasPrefix(*text, "#") {
		*skipLine = true
		*skipMultipleLines = false
	} else if *text == "" {
		*skipLine = true
		*skipMultipleLines = false
	} else if strings.HasPrefix(*text, "=begin") { // multiline comment begins
		*skipLine = true
		*skipMultipleLines = true
	} else {
		*skipLine = false
		*skipMultipleLines = false
	}
}

func parsePythonFile(text *string, effectiveCount *int, skipLine *bool, skipMultipleLines *bool) {
	if *skipMultipleLines {
		*skipLine = true
		if strings.Contains(*text, "\"\"\"") {
			*skipMultipleLines = false
		} else {
			*skipMultipleLines = true
		}
		return
	}

	if strings.HasPrefix(*text, "#") {
		*skipLine = true
		*skipMultipleLines = false
	} else if *text == "" {
		*skipLine = true
		*skipMultipleLines = false
	} else if *text == "\"\"\"" {
		*skipLine = true
		if *skipMultipleLines {
			*skipMultipleLines = false
		} else {
			*skipMultipleLines = true
		}
	} else if strings.HasPrefix(*text, "\"\"\"") { // docstring on one line
		*skipLine = true
		if strings.HasSuffix(*text, "\"\"\"") { // closes on the same line
			*skipMultipleLines = false
		} else {
			*skipMultipleLines = true
		}
	} else {
		*skipLine = false
		*skipMultipleLines = false
	}
}

func parsePHPFile(text *string, effectiveCount *int, skipLine *bool, skipMultipleLines *bool) {
	// PHP comments are similar to C style comments + the "#" as comment
	parseCFile(text, effectiveCount, skipLine, skipMultipleLines)
	// also include "#" as comment
	if strings.HasPrefix(*text, "#") {
		*skipLine = true
		*skipMultipleLines = false
	}
}

// parses every file that has C like comments and is in the list of valid files
func parseCFile(text *string, effectiveCount *int, skipLine *bool, skipMultipleLines *bool) {
	if strings.HasPrefix(*text, "//") { // Handle single line comments
		*skipLine = true
		*skipMultipleLines = false
	} else if *text == "" { // check for blank lines
		*skipLine = true
		*skipMultipleLines = false
	} else if strings.HasPrefix(*text, "/*") { // check for multiline comments
		*skipLine = true
		if strings.Contains(*text, "*/") && !strings.Contains(*text, "\"*/") &&
			!strings.Contains(*text, "*/\"") && !strings.Contains(*text, "'*/") &&
			!strings.Contains(*text, "*/'") { // check for close on the same line
			if !strings.HasSuffix(*text, "/") {
				*skipLine = false
			}
			*skipMultipleLines = false
		} else {
			*skipMultipleLines = true // switch to skip following lines until comment closes
		}
	} else if strings.Contains(*text, "/*") && !strings.Contains(*text, "\"/*") &&
		!strings.Contains(*text, "/*\"") && !strings.Contains(*text, "'/*") &&
		!strings.Contains(*text, "*/'") { // multiline comment that starts after some code in that line
		*skipLine = false
		if strings.HasSuffix(*text, "*/") {
			*skipMultipleLines = false
		} else {
			*skipMultipleLines = true
		}
	} else {
		*skipLine = false
	}

	// check to see if we are still in a comment
	if *skipMultipleLines {
		if strings.Contains(*text, "*/") && !strings.Contains(*text, "\"*/") &&
			!strings.Contains(*text, "*/\"") && !strings.Contains(*text, "'*/") &&
			!strings.Contains(*text, "*/'") {
			*skipMultipleLines = false
			if !strings.HasSuffix(*text, "*/") {
				*skipLine = false
			} else {
				*skipLine = true
			}
		} else {
			*skipMultipleLines = true
			if !strings.HasPrefix(*text, "/*") && strings.Contains(*text, "/*") &&
				!strings.Contains(*text, "\"/*") && !strings.Contains(*text, "/*\"") &&
				!strings.Contains(*text, "'/*") && !strings.Contains(*text, "*/'") {
				*skipLine = false
			} else {
				*skipLine = true
			}
		}
	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filenames, err := getFilenames()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesAndLines := make(map[string]int)
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		countLines(file, filesAndLines)
		defer file.Close()
	}

	printLineNumbers(filesAndLines)
}

func getFilenames() ([]string, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("a filename or folder ('.' for the current directory) needs to be entered")
	}

	subDirsToSkip := [...]string{"venv", "node_modules", "build", "bin", "out"}

	var filenames []string
	// Entering a dot as argument means scan whole directory, including subdirectories
	if os.Args[1] == "." {
		rootDir := "."
		filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			validFile := strings.HasSuffix(info.Name(), ".go") ||
				strings.HasSuffix(info.Name(), ".c") || strings.HasSuffix(info.Name(), ".cpp") ||
				strings.HasSuffix(info.Name(), ".cc") || strings.HasSuffix(info.Name(), ".cs") ||
				strings.HasSuffix(info.Name(), ".h") || strings.HasSuffix(info.Name(), ".hpp") ||
				strings.HasSuffix(info.Name(), ".py") || strings.HasSuffix(info.Name(), ".java") ||
				strings.HasSuffix(info.Name(), ".js") || strings.HasSuffix(info.Name(), ".jsx") ||
				strings.HasSuffix(info.Name(), ".ts") || strings.HasSuffix(info.Name(), ".php") ||
				strings.HasSuffix(info.Name(), ".rs") || strings.HasSuffix(info.Name(), ".rb") ||
				strings.HasSuffix(info.Name(), ".R")

			// Skip unwanted directories (build, hidden etc.)
			hiddenDir := info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != rootDir
			skipDir := false
			for _, dir := range subDirsToSkip {
				if info.Name() == dir {
					skipDir = true
					break
				}
				skipDir = false
			}
			if info.IsDir() && (skipDir || hiddenDir) {
				fmt.Printf("Skipping unwanted dir: %+v \n", path)
				return filepath.SkipDir
			}

			// Only append valid source files to the list
			if !info.IsDir() && validFile {
				filenames = append(filenames, path)
			}
			return nil
		})
	} else {
		// @TODO: Make folder names possible too, e.g. "loc src" searches the src directory
		filenames = append(filenames, os.Args[1:]...)
	}

	return filenames, nil
}

func countLines(file *os.File, fileMap map[string]int) {
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("%s in file: %v\n", err, file.Name())
	}
	fileMap[file.Name()] = lineCount
}

func printLineNumbers(filesAndLines map[string]int) {
	totalCount := 0
	for filename, count := range filesAndLines {
		fmt.Printf("%v: %v\n", filename, count)
		totalCount += count
	}

	fmt.Println("The total line count is", totalCount)
}

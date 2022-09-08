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
	for _, rootDir := range os.Args[1:] {
		filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			validFile := strings.HasSuffix(d.Name(), ".go") ||
				strings.HasSuffix(d.Name(), ".c") || strings.HasSuffix(d.Name(), ".cpp") ||
				strings.HasSuffix(d.Name(), ".cc") || strings.HasSuffix(d.Name(), ".cs") ||
				strings.HasSuffix(d.Name(), ".h") || strings.HasSuffix(d.Name(), ".hpp") ||
				strings.HasSuffix(d.Name(), ".py") || strings.HasSuffix(d.Name(), ".java") ||
				strings.HasSuffix(d.Name(), ".js") || strings.HasSuffix(d.Name(), ".jsx") ||
				strings.HasSuffix(d.Name(), ".ts") || strings.HasSuffix(d.Name(), ".php") ||
				strings.HasSuffix(d.Name(), ".rs") || strings.HasSuffix(d.Name(), ".rb") ||
				strings.HasSuffix(d.Name(), ".R")

			// Skip unwanted directories (build, hidden etc.)
			hiddenDir := d.IsDir() && strings.HasPrefix(d.Name(), ".") && d.Name() != rootDir
			skipDir := false
			for _, dir := range subDirsToSkip {
				if d.Name() == dir {
					skipDir = true
					break
				}
				skipDir = false
			}
			if d.IsDir() && (skipDir || hiddenDir) {
				fmt.Printf("Skipping unwanted dir: %+v \n", path)
				return filepath.SkipDir
			}

			// Only append valid source files to the list
			if !d.IsDir() && validFile {
				filenames = append(filenames, path)
			}
			return nil
		})
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

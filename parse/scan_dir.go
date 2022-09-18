package parse

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"loc/options"
	"path/filepath"
	"strings"
)

func GetFilenames() ([]string, error) {
	if flag.NArg() < 1 {
		return nil, errors.New("a filename or folder ('.' for the current directory) needs to be entered")
	}

	subDirsToSkip := [...]string{"venv", "node_modules", "build", "bin", "out"}

	var filenames []string
	// Entering a dot as argument means scan whole directory, including subdirectories
	for _, rootDir := range flag.Args() {
		filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("Could not access file %q\n", path)
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
				if options.Verbose {
					fmt.Printf("Skipping unwanted dir: %+v \n", path)
				}
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

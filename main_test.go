package main

import (
	"os"
	"testing"
)

func TestCountLines(t *testing.T) {
	// Only languages with different comment styles need to be tested (e.g. Go has the same comments as C/C++ etc.)
	files := [...]string{
		"testfiles/test.go", "testfiles/test.rs", "testfiles/test.py",
		"testfiles/test.php", "testfiles/test.rb", "testfiles/test.R",
		"testfiles/test_comment_before_after.go", "testfiles/test_advanced.go",
	}
	os.Args = []string{"cmd"}
	for i := 0; i < len(files); i++ {
		os.Args = append(os.Args, files[i])
	}
	filenames, err := getFilenames()
	if err != nil {
		t.Errorf("Test FAILED: Could not open file")
	}

	filesAndLines := make(map[string][2]int)
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			t.Errorf("Test FAILED: Could not open file %v", file.Name())
		}
		defer file.Close()
		countLines(file, filesAndLines)
	}

	for file := range filesAndLines {
		var expected [2]int
		switch file {
		case "testfiles/test.go":
			expected = [2]int{26, 12}
		case "testfiles/test_comment_before_after.go":
			expected = [2]int{5, 3}
		case "testfiles/test_advanced.go":
			expected = [2]int{16, 4}
		case "testfiles/test.rs":
			expected = [2]int{61, 9}
		case "testfiles/test.py":
			expected = [2]int{25, 9}
		case "testfiles/test.php":
			expected = [2]int{15, 5}
		case "testfiles/test.rb":
			expected = [2]int{10, 2}
		case "testfiles/test.R":
			expected = [2]int{6, 1}
		}

		if filesAndLines[file] == expected {
			t.Logf("%v SUCCESS: Expected %v, got %v", file, expected, filesAndLines[file])
		} else {
			t.Errorf("%v FAILED: Expected %v, got %v", file, expected, filesAndLines[file])
		}
	}
}

package main

import (
	"os"
	"testing"
)

func TestCountLines(t *testing.T) {
	filesAndLines := make(map[string][2]int)
	file, err := os.Open("testfiles/test.go")
	if err != nil {
		t.Errorf("Test FAILED: Could not open file %v", file)
	}
	defer file.Close()

	countLines(file, filesAndLines)
	expected := [2]int{23, 9}
	if filesAndLines["testfiles/test.go"] == expected {
		t.Logf("Test SUCCESS: Total and effective lines counted correctly")
	} else {
		t.Errorf("Test FAILED: Expected %v, got %v", expected, filesAndLines["testfiles/test.go"])
	}
}

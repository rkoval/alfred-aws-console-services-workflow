package util

import "testing"

func TestGetCurrentFilename(t *testing.T) {
	currentFilename := GetCurrentFilename()
	if currentFilename != "tools_test" {
		t.Errorf("filename was not correct: %s", currentFilename)
	}
}

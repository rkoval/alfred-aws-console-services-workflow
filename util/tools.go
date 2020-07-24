package util

import (
	"path/filepath"
	"runtime"
	"strings"
)

func GetCurrentFilename() string {
	_, current_file, _, _ := runtime.Caller(1)
	baseFile := filepath.Base(current_file)
	return strings.TrimSuffix(baseFile, filepath.Ext(baseFile))
}

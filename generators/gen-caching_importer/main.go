package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	importLines := getImportLines()

	filename := "caching/gen-caching.go"
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content[:]), "\n")

	extraImports := []string{
		"github.com/aws/aws-sdk-go-v2/aws",
		"github.com/aws/smithy-go",
	}
	for i, line := range lines {
		if strings.Contains(line, "\"github.com/deanishe/awgo\"") {
			lines = append(lines[:i+1+len(importLines)+len(extraImports)], lines[i+1:]...)
			for j, importLine := range importLines {
				lines[i+1+j] = importLine
			}
			for k, extraImportLine := range extraImports {
				lines[i+1+len(importLines)+k] = "\t\"" + extraImportLine + "\""
			}
			break
		}
	}

	err = os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0600)
	if err != nil {
		panic(err)
	}
}

func getImportLines() []string {
	cachingFileContents, err := os.ReadFile("caching/caching.go")
	if err != nil {
		panic(err)
	}
	entitiesRawString := regexp.MustCompile("//go:generate genny .*\"Entity=(.+)\"")
	matches := entitiesRawString.FindStringSubmatch(string(cachingFileContents[:]))
	entities := strings.Split(matches[1], ",")

	importLines := []string{}
	seenPackages := map[string]bool{}
	for _, entity := range entities {
		pkg := strings.Split(entity, ".")[0]
		if !seenPackages[pkg] {
			importLines = append(importLines, fmt.Sprintf("\t%s \"github.com/aws/aws-sdk-go-v2/service/%s/types\"", pkg, pkg))
		}
		seenPackages[pkg] = true
	}

	return importLines
}

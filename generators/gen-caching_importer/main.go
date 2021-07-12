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

	for i, line := range lines {
		if strings.Contains(line, "\"github.com/aws/aws-sdk-go-v2/aws\"") {
			lines = append(lines[:i+1+len(importLines)], lines[i+1:]...) // i < len(a)
			for j, importLine := range importLines {
				lines[i+1+j] = importLine
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
	for _, entity := range entities {
		pkg := strings.Split(entity, ".")[0]
		importLines = append(importLines, fmt.Sprintf("\t%s \"github.com/aws/aws-sdk-go-v2/service/%s/types\"", pkg, pkg))
	}

	return importLines
}

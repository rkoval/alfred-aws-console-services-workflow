package main

import (
	"io/ioutil"
	"sort"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("../searchers/searchers_by_service_id.go")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content[:]), "\n")
	var startIndexOfInstances int
	var endIndexOfInstances int
	var startIndexOfMapEntries int
	var endIndexOfMapEntries int
	for i, line := range lines {
		lineHasInstance := strings.HasSuffix(line, "Searcher{}")
		if lineHasInstance && startIndexOfInstances == 0 {
			startIndexOfInstances = i
		} else if !lineHasInstance && startIndexOfInstances != 0 && endIndexOfInstances == 0 {
			endIndexOfInstances = i
			sort.Strings(lines[startIndexOfInstances:endIndexOfInstances])
		}

		lineHasMapEntry := strings.HasSuffix(line, ",")
		if lineHasMapEntry && startIndexOfMapEntries == 0 {
			startIndexOfMapEntries = i
		} else if startIndexOfMapEntries != 0 && !lineHasMapEntry && endIndexOfMapEntries == 0 {
			endIndexOfMapEntries = i
			sort.Strings(lines[startIndexOfMapEntries:endIndexOfMapEntries])
			break
		}
	}

	err = ioutil.WriteFile("../searchers/searchers_by_service_id.go", []byte(strings.Join(lines, "\n")), 0600)
	if err != nil {
		panic(err)
	}
}

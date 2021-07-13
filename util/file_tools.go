package util

import (
	"os"
	"regexp"
	"strings"
	"text/template"
)

func WriteTemplateToFile(templateName, templateString, fileName string, data interface{}) {
	t, err := template.New(templateName).Parse(templateString)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = t.Execute(file, data)
	if err != nil {
		panic(err)
	}
}

func ModifyFileWithRegexReplace(filename string, regex *regexp.Regexp, replacement string, ignoreIfContains string) string {
	c, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	content := string(c)

	if ignoreIfContains != "" && strings.Contains(content, ignoreIfContains) {
		return ""
	}

	replacedContent := regex.ReplaceAllString(content, replacement)
	err = os.WriteFile(filename, []byte(replacedContent), 0600)
	if err != nil {
		panic(err)
	}
	return replacedContent
}

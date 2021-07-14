package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func appendToGenny(searcherNamer SearcherNamer) {
	regex := regexp.MustCompile(`(go:generate genny.*)(")`)
	typeName := searcherNamer.OperationDefinition.Package + "." + searcherNamer.OperationDefinition.Item
	replacement := fmt.Sprintf("$1,%s$2", typeName)
	util.ModifyFileWithRegexReplace("caching/caching.go", regex, replacement, typeName)
}

func appendToSearchers(searcherNamer SearcherNamer) {
	regex := regexp.MustCompile(`(\nvar SearchersByServiceId)`)
	structInitializer := "&" + searcherNamer.StructName + "{}"
	replacement := fmt.Sprintf("var %s = %s\n$1", searcherNamer.StructInstanceName, structInitializer)
	replacedContent := util.ModifyFileWithRegexReplace("searchers/searchers_by_service_id.go", regex, replacement, structInitializer)
	if replacedContent == "" {
		return
	}

	regex = regexp.MustCompile(`(,\n)(})`)
	replacement = fmt.Sprintf("$1\t\"%s\": %s$1$2", searcherNamer.NameSnakeCasePlural, searcherNamer.StructInstanceName)
	replacedContent = util.ModifyFileWithRegexReplace("searchers/searchers_by_service_id.go", regex, replacement, "")

	if !strings.Contains(replacedContent, fmt.Sprintf("\"%s\"", searcherNamer.ServiceLower)) {
		// append root service if this is the first one we're populating
		regex = regexp.MustCompile(`(,\n)(})`)
		replacement = fmt.Sprintf("$1\t\"%s\": %s$1$2", searcherNamer.ServiceLower, searcherNamer.StructInstanceName)
		util.ModifyFileWithRegexReplace("searchers/searchers_by_service_id.go", regex, replacement, "")
	}
}

func appendToWorkflowTest(searcherNamer SearcherNamer) {
	filename := "workflow/workflow_test.go"
	regex := regexp.MustCompile(`(\t},)(\n})`)

	addWorkflowTestCaseString := func(query string) string {
		testCaseString := fmt.Sprintf(`	{
		query:       "%s",
		fixtureName: "../searchers/%s_test", // reuse test fixture from this other test
	},`, query, searcherNamer.NameSnakeCasePlural)
		replacement := fmt.Sprintf("$1\n%s$2", testCaseString)
		return util.ModifyFileWithRegexReplace(filename, regex, replacement, "\""+query+"\"")
	}

	addWorkflowTestCaseString(searcherNamer.ServiceLower)
	addWorkflowTestCaseString(searcherNamer.ServiceLower + " ")
	addWorkflowTestCaseString(searcherNamer.ServiceLower + " " + searcherNamer.EntityLowerPlural)
	addWorkflowTestCaseString(searcherNamer.ServiceLower + " " + searcherNamer.EntityLowerPlural + " ")
}

func writeSearcherFile(searcherNamer SearcherNamer) {
	templateString := `package searchers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/{{ .OperationDefinition.Package }}"
	"github.com/aws/aws-sdk-go-v2/service/{{ .OperationDefinition.Package }}/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type {{ .StructName }} struct{}

func (s {{ .StructName }}) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.Load{{ .OperationDefinition.PackageTitle }}{{ .OperationDefinition.Item }}ArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
	}
	return nil
}

func (s {{ .StructName }}) fetch(cfg aws.Config) ([]types.{{ .OperationDefinition.Item }}, error) {
	client := {{ .OperationDefinition.Package }}.NewFromConfig(cfg)

	entities := []types.{{ .OperationDefinition.Item }}{}
	{{if .OperationDefinition.PageInputToken }}pageToken := ""
	for { {{ end }}
		params := &{{ .OperationDefinition.Package }}.{{ .OperationDefinition.FunctionInput }}{
			{{if .OperationDefinition.PageSize }}{{ .OperationDefinition.PageSize }}: aws.Int32(1000),{{ end }}
		}
		{{if .OperationDefinition.PageInputToken }}if pageToken != "" {
			params.{{ .OperationDefinition.PageInputToken }} = &pageToken
		}{{ end }}
		resp, err := client.{{ .OperationDefinition.FunctionName }}(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.{{ .OperationDefinition.Items }}...)

		{{if .OperationDefinition.PageOutputToken }}if resp.{{ .OperationDefinition.PageOutputToken }} != nil {
			pageToken = *resp.{{ .OperationDefinition.PageOutputToken }}
		} else {
			break
		}{{ end }}
	{{if .OperationDefinition.PageInputToken }} }{{ end }}

	return entities, nil
}

func (s {{ .StructName }}) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, entity types.{{ .OperationDefinition.Item }}) {
	title := entity.TODO
	subtitle := ""

	path := fmt.Sprintf("/{{ .ServiceLower }}/{{ .EntityLowerPlural }}/?region=%s", config.Region)
	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, config.Region)).
		Icon(awsworkflow.GetImageIcon("{{ .ServiceLower }}")).
		Valid(true)
}`

	util.WriteTemplateToFile("searcher_file", templateString, fmt.Sprintf("searchers/%s.go", searcherNamer.NameSnakeCasePlural), searcherNamer)
}

func writeSearcherTestFile(searcherNamer SearcherNamer) {
	templateString := `package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Test{{ .StructName }}(t *testing.T) {
	TestSearcher(t, {{ .StructName }}{}, util.GetCurrentFilename())
}`

	util.WriteTemplateToFile("searcher_test_file", templateString, fmt.Sprintf("searchers/%s_test.go", searcherNamer.NameSnakeCasePlural), searcherNamer)
}

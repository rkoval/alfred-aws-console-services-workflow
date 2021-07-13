package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/iancoleman/strcase"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
)

func init() {
	awsServices := parsers.ParseConsoleServicesYml("./console-services.yml")
	for _, awsService := range awsServices {
		if awsService.ShortName != "" {
			strcase.ConfigureAcronym(awsService.ShortName, awsService.ShortName)
		}
	}
	flag.Parse()
}

type OperationDefinition struct {
	Package         string
	PackageTitle    string
	FunctionName    string
	FunctionInput   string
	Item            string
	Items           string
	PageInputToken  string
	PageOutputToken string
	PageSize        string
}

type SearcherNamer struct {
	ServiceTitle       string
	ServiceLower       string
	EntityTitle        string
	EntityLower        string
	EntityLowerPlural   string
	Name               string
	NameLower          string
	NameCamelCase      string
	NameSnakeCase      string
	NameSnakeCasePlural string
	StructName         string
	StructInstanceName string
	OperationDefinition
}

func NewSearcherNamer(service, entity string, operationDefinition OperationDefinition) SearcherNamer {
	if "s" == entity[len(entity)-1:] {
		log.Fatalf("Entity should be singular for casing to work properly")
	}

	serviceTitle := strings.Title(service)
	entityTitle := strings.Title(entity)
	serviceLower := strings.ToLower(service)
	name := serviceTitle + entityTitle

	return SearcherNamer{
		ServiceTitle:        serviceTitle,
		ServiceLower:        serviceLower,
		EntityTitle:         entityTitle,
		EntityLower:         strings.ToLower(entity),
		EntityLowerPlural:   strings.ToLower(entity) + "s", // TODO make this proper english
		Name:                name,
		NameLower:           strings.ToLower(name),
		NameCamelCase:       strcase.ToCamel(name),
		NameSnakeCase:       strcase.ToSnake(name),
		NameSnakeCasePlural: strcase.ToSnake(name) + "s", // TODO make this proper english
		StructName:          name + "Searcher",
		StructInstanceName:  serviceLower + entityTitle + "Searcher",
		OperationDefinition: operationDefinition,
	}
}

func main() {
	args := flag.Args()
	if len(args) < 3 {
		usage()
	}

	operation := args[2]
	pkg, functionName := parseOperation(operation)
	goGetPkg(pkg)

	operationDefinition := getOperationDefinition(operation, pkg, functionName)
	searcherNamer := NewSearcherNamer(args[0], args[1], operationDefinition)

	appendToGenny(searcherNamer)
	appendToSearchers(searcherNamer)
	writeSearcherFile(searcherNamer)
	writeSearcherTestFile(searcherNamer)
}

func goGetPkg(pkg string) {
	cmd := exec.Command("go", "get", "github.com/aws/"+aws.SDKName+"/service/"+pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func parseOperation(operation string) (string, string) {
	operationNameRegex := regexp.MustCompile("com.amazonaws.([a-z]+)#([a-zA-Z]+)")
	matches := operationNameRegex.FindStringSubmatch(operation)
	if len(matches) != 3 {
		log.Fatalln("operation argument must have the form \"com.amazonaws.pkg#FunctionName\"")
	}
	return matches[1], matches[2]
}

func getOperationDefinition(operation, pkg, functionName string) OperationDefinition {
	gopath, exists := os.LookupEnv("GOPATH")
	if !exists {
		userHome, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		gopath += userHome + "/go"
	}

	globPath := gopath + "/pkg/mod/github.com/aws/" + aws.SDKName + "@v" + aws.SDKVersion + "/codegen/sdk-codegen/aws-models/" + pkg + ".*.json"
	matches, err := filepath.Glob(globPath)
	if err != nil {
		panic(err)
	}
	if len(matches) <= 0 {
		panic(errors.New("Unable to find a file with glob \"" + globPath + "\""))
	} else if len(matches) >= 2 {
		panic(errors.New("More than one file with glob \"" + globPath + "\""))
	}
	filename := matches[0]

	apiJsonRaw, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var j interface{}
	err = json.Unmarshal(apiJsonRaw, &j)
	if err != nil {
		panic(err)
	}

	definition := getJsonPath(j, "shapes", operation).(map[string]interface{})

	_, functionInput := parseOperation(getJsonPath(definition, "input", "target").(string))
	functionOutputShape := getJsonPath(definition, "output", "target").(string)
	paginated := getJsonPath(definition, "traits", "smithy.api#paginated").(map[string]interface{})
	items := paginated["items"].(string)

	functionOutputItemsShape := getJsonPath(j, "shapes", functionOutputShape, "members", items, "target").(string)
	_, item := parseOperation(getJsonPath(j, "shapes", functionOutputItemsShape, "member", "target").(string))

	operationDefinition := OperationDefinition{
		Package:       pkg,
		PackageTitle:  strings.Title(pkg),
		FunctionName:  functionName,
		FunctionInput: functionInput,
		Item:          item,
		Items:         items,
	}

	pageInputToken := paginated["inputToken"]
	if pageInputToken != nil {
		operationDefinition.PageInputToken = pageInputToken.(string)
	}
	pageOutputToken := paginated["outputToken"]
	if pageOutputToken != nil {
		operationDefinition.PageOutputToken = pageOutputToken.(string)
	}
	pageSize := paginated["pageSize"]
	if pageSize != nil {
		operationDefinition.PageSize = pageSize.(string)
	}

	return operationDefinition
}

func getJsonPath(json interface{}, keys ...string) interface{} {
	value := json
	for _, key := range keys {
		value = value.(map[string]interface{})[key]
	}

	return value
}

func usage() {
	flag.Usage()
	fmt.Println("go run searcher.go Service Entity com.amazonaws.package#functionName")
	os.Exit(1)
}

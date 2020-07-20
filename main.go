package main

import (
	"io/ioutil"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"github.com/rkoval/alfred-aws-console-services-workflow/filters"
	"gopkg.in/yaml.v2"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func parseYaml() []core.AwsService {
	awsServices := []core.AwsService{}
	yamlFile, err := ioutil.ReadFile("console-services.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &awsServices)
	if err != nil {
		log.Fatal(err)
	}
	return awsServices
}

func run() {
	var query string
	if args := wf.Args(); len(args) > 0 {
		query = strings.TrimSpace(args[0])
	}

	awsServices := parseYaml()

	awsServicesById := make(map[string]*core.AwsService)
	for i, awsService := range awsServices {
		awsServicesById[awsService.Id] = &awsServices[i]
	}

	// TODO add better lexing here to route filters

	splitQuery := strings.Split(query, " ")
	if len(splitQuery) <= 1 || awsServicesById[splitQuery[0]] == nil {
		filters.Services(wf, awsServices, query)
	} else {
		awsService := awsServicesById[splitQuery[0]]
		filters.ServiceSections(wf, *awsService, strings.Join(splitQuery[1:], " "))
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}

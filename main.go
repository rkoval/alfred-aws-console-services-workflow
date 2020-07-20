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
	if len(splitQuery) <= 1 || awsServicesById[splitQuery[0]] == nil || len(awsServicesById[splitQuery[0]].Sections) <= 0 {
		filters.Services(wf, awsServices, query)
	} else {
		id := splitQuery[0]
		log.Printf("filtering on sections for %s", id)
		awsService := awsServicesById[id]
		filters.ServiceSections(wf, *awsService, strings.Join(splitQuery[1:], " "))
	}

	if query != "" {
		res := wf.Filter(query)

		log.Printf("%d results match %q", len(res), query)

		for i, r := range res {
			log.Printf("%02d. score=%0.1f sortkey=%s", i+1, r.Score, wf.Feedback.Keywords(i))
		}
	}

	wf.WarnEmpty("No matching services found", "Try a different query?")

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}

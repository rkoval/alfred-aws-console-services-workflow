package main

import (
	"io/ioutil"
	"log"

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
		query = args[0]
	}

	awsServices := parseYaml()

	// TODO add lexing here to route filters

	filters.Services(wf, awsServices, query)

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}

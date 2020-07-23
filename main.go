package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"gopkg.in/yaml.v2"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func readConsoleServicesYml() []core.AwsService {
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

func Run(wf *aw.Workflow, query string, transport http.RoundTripper) {
	awsServices := readConsoleServicesYml()
	query, err := parsers.ParseQueryAndPopulateItems(wf, awsServices, query, transport)

	if err != nil {
		wf.FatalError(err)
	} else if query != "" {
		log.Printf("filtering with query %s", query)
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
	var query string
	args := wf.Args()
	log.Printf("running workflow with %d arg(s): %v", len(args), args)
	if len(args) > 0 {
		query = strings.TrimLeft(args[0], " ")
	}

	wf.Run(func() {
		Run(wf, query, nil)
	})
}

package main

import (
	"flag"
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
var forceFetch bool
var query string

func init() {
	flag.BoolVar(&forceFetch, "fetch", false, "force fetch via AWS instead of cache")
	flag.StringVar(&query, "query", "", "query to use")
	flag.Parse()
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
	query = parsers.ParseQueryAndPopulateItems(wf, awsServices, query, transport, forceFetch)

	if query != "" {
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
	wf.Run(func() {
		log.Printf("running workflow with query: `%s`", query)
		query = strings.TrimLeft(query, " ")

		Run(wf, query, nil)
	})
}

package workflow

import (
	"io/ioutil"
	"log"
	"net/http"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"gopkg.in/yaml.v2"
)

func readConsoleServicesYml(ymlPath string) []core.AwsService {
	awsServices := []core.AwsService{}
	yamlFile, err := ioutil.ReadFile(ymlPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &awsServices)
	if err != nil {
		log.Fatal(err)
	}
	return awsServices
}

func Run(wf *aw.Workflow, query string, transport http.RoundTripper, forceFetch bool, ymlPath string) {
	awsServices := readConsoleServicesYml(ymlPath)
	query = ParseQueryAndPopulateItems(wf, awsServices, query, transport, forceFetch)

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

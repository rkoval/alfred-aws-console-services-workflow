package main

import (
	"io/ioutil"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
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

func populateItems(awsServices []core.AwsService, query string) (string, error) {
	awsServicesById := make(map[string]*core.AwsService)
	for i, awsService := range awsServices {
		awsServicesById[awsService.Id] = &awsServices[i]
	}

	// TODO add better lexing here to route searchers

	splitQuery := strings.Split(query, " ")
	if len(splitQuery) > 1 && awsServicesById[splitQuery[0]] != nil {
		id := splitQuery[0]
		query = strings.Join(splitQuery[1:], " ")
		awsService := awsServicesById[id]
		searcher := searchers.SearchersByServiceId[id]
		if strings.HasPrefix(query, "$") && searcher != nil {
			query = query[1:]
			log.Printf("using searcher associated with %s", id)
			err := searcher(wf, query, nil)
			if err != nil {
				return "", err
			}
			return "", nil
		} else if len(awsServicesById[id].Sections) > 0 {
			sections := awsServicesById[id].Sections
			sectionsById := make(map[string]*core.AwsServiceSection)
			for i, section := range sections {
				sectionsById[section.Id] = &sections[i]
			}
			if len(splitQuery) > 2 && sectionsById[splitQuery[1]] != nil {
				sectionId := splitQuery[1]
				query = strings.Join(splitQuery[2:], " ")
				id = id + "_" + sectionId
				searcher := searchers.SearchersByServiceId[id]
				if searcher != nil {
					log.Printf("using searcher associated with %s", id)
					err := searcher(wf, query, nil)
					if err != nil {
						return "", err
					}
					return "", nil
				}
			}
			log.Printf("filtering on sections for %s", id)
			query = strings.TrimSpace(strings.Join(splitQuery[1:], " "))
			searchers.ServiceSections(wf, *awsService, query)
			return query, nil
		}
	}

	searchers.Services(wf, awsServices, query)
	return query, nil
}

func run() {
	var query string
	if args := wf.Args(); len(args) > 0 {
		query = strings.TrimLeft(args[0], " ")
	}

	awsServices := parseYaml()

	query, err := populateItems(awsServices, query)

	if err != nil {
		log.Printf("error: %v", err)
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
	wf.Run(run)
}

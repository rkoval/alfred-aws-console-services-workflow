package parsers

import (
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
)

func ParseQueryAndPopulateItems(wf *aw.Workflow, awsServices []core.AwsService) (string, error) {
	var query string
	if args := wf.Args(); len(args) > 0 {
		query = strings.TrimLeft(args[0], " ")
	}

	// TODO break apart this function
	// TODO add better lexing here to route searchers

	splitQuery := strings.Split(query, " ")
	if len(splitQuery) > 1 {
		id := splitQuery[0]
		var awsService *core.AwsService
		for i := range awsServices {
			if awsServices[i].Id == id {
				awsService = &awsServices[i]
				break
			}
		}

		if awsService != nil {
			searcher := searchers.SearchersByServiceId[id]
			if strings.HasPrefix(query, "$") && searcher != nil {
				query = query[1:]
				log.Printf("using searcher associated with %s", id)
				err := searcher(wf, query, nil)
				if err != nil {
					return "", err
				}
				return "", nil
			} else if len(awsService.Sections) > 0 {
				sections := awsService.Sections
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
				searchers.SearchServiceSections(wf, *awsService)
				return query, nil
			} else {
				// no sections defined, just fake out home as a sub-section
				awsService.Sections = []core.AwsServiceSection{
					{
						Id:   "home",
						Name: "Home",
						Url:  awsService.Url,
					},
				}

				searchers.SearchServiceSections(wf, *awsService)
				return "", nil
			}
		}
	}

	searchers.SearchServices(wf, awsServices)
	return query, nil
}

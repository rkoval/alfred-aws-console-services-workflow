package parsers

import (
	"log"
	"net/http"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
)

func ParseQueryAndPopulateItems(wf *aw.Workflow, awsServices []core.AwsService, query string, transport http.RoundTripper) (string, error) {
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
			query = strings.Join(splitQuery[1:], " ")
			searcher := searchers.SearchersByServiceId[id]
			if strings.HasPrefix(query, "$") && searcher != nil {
				query = query[1:]
				log.Printf("using searcher associated with %s", id)
				err := searcher(wf, query, transport)
				if err != nil {
					return "", err
				}
				return "", nil
			} else {
				// prepend the home to the sub-service list so that it's still accessible
				awsServiceHome := *awsService
				awsServiceHome.Id = "home"
				awsService.SubServices = append(
					[]core.AwsService{
						awsServiceHome,
					},
					awsService.SubServices...,
				)

				if len(awsService.SubServices) > 1 {
					splitQuery = strings.Split(query, " ")
					if len(splitQuery) > 1 {
						subServiceId := splitQuery[0]
						var subService *core.AwsService
						for i := range awsService.SubServices {
							if awsService.SubServices[i].Id == subServiceId {
								subService = &awsService.SubServices[i]
								break
							}
						}
						if subService != nil {
							query = strings.Join(splitQuery[1:], " ")
							id = id + "_" + subServiceId
							log.Println("id", id)
							searcher := searchers.SearchersByServiceId[id]
							if searcher != nil {
								log.Printf("using searcher associated with %s", id)
								err := searcher(wf, query, transport)
								if err != nil {
									return "", err
								}
								return "", nil
							}
						}
					}
				}
				log.Printf("filtering on subServices for %s", id)
				query = strings.TrimSpace(strings.Join(splitQuery, " "))
				searchers.SearchSubServices(wf, *awsService)
				return query, nil
			}
		}
	}

	searchers.SearchServices(wf, awsServices)
	return query, nil
}

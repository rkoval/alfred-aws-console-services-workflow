package workflow

import (
	"log"
	"net/http"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

func ParseQueryAndPopulateItems(wf *aw.Workflow, awsServices []core.AwsService, query string, transport http.RoundTripper, forceFetch bool) string {
	// TODO break apart this function
	// TODO add better lexing here to route searchers

	fullQuery := query

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
			searcher := SearchersByServiceId[id]
			if strings.HasPrefix(query, "$") && searcher != nil {
				query = query[1:]
				log.Printf("using searcher associated with %s", id)
				err := searcher(wf, query, transport, forceFetch, fullQuery)
				if err != nil {
					wf.FatalError(err)
				}
				return query
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
							searcher := SearchersByServiceId[id]
							if searcher != nil {
								log.Printf("using searcher associated with %s", id)
								err := searcher(wf, query, transport, forceFetch, fullQuery)
								if err != nil {
									wf.FatalError(err)
								}
								return query
							}
						}
					}
				}
				log.Printf("filtering on subServices for %s", id)
				query = strings.TrimSpace(strings.Join(splitQuery, " "))
				SearchSubServices(wf, *awsService)
				return query
			}
		}
	}

	SearchServices(wf, awsServices)
	return query
}

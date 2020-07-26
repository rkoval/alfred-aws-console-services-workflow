package parsers

import (
	"os"
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchtypes"
)

func ParseQuery(awsServices []awsworkflow.AwsService, query string) (string, searchtypes.SearchType, *awsworkflow.AwsService, bool) {
	// TODO add better lexing here to route populators

	searchAlias := os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_SEARCH_ALIAS")
	if searchAlias == "" {
		searchAlias = ","
	}

	splitQuery := strings.Split(query, " ")
	if len(splitQuery) > 1 {
		id := splitQuery[0]
		var awsService *awsworkflow.AwsService
		for i := range awsServices {
			if awsServices[i].Id == id {
				awsService = &awsServices[i]
				break
			}
		}

		if awsService != nil {
			query = strings.Join(splitQuery[1:], " ")
			searchType := searchtypes.SearchTypesByServiceId[id]
			if strings.HasPrefix(query, searchAlias) && searchType != 0 {
				query = query[len(searchAlias):]
				return query, searchType, awsService, false
			} else if len(awsService.SubServices) > 0 {
				// prepend the home to the sub-service list so that it's still accessible
				if awsService.HomeID == "" {
					awsServiceHome := *awsService
					awsServiceHome.Id = "home"
					awsService.SubServices = append(
						[]awsworkflow.AwsService{
							awsServiceHome,
						},
						awsService.SubServices...,
					)
				}

				if len(awsService.SubServices) > 1 {
					if query == "OPEN_ALL" {
						return query, searchtypes.SubServices, awsService, true
					}
					splitQuery = strings.Split(query, " ")
					if len(splitQuery) > 1 {
						subServiceId := splitQuery[0]
						var subService *awsworkflow.AwsService
						for i := range awsService.SubServices {
							if awsService.SubServices[i].Id == subServiceId {
								subService = &awsService.SubServices[i]
								break
							}
						}
						if subService != nil {
							query = strings.Join(splitQuery[1:], " ")
							id = id + "_" + subServiceId
							searchType := searchtypes.SearchTypesByServiceId[id]
							if searchType != 0 {
								return query, searchType, subService, false
							}
						}
					}
				}

				query = strings.TrimSpace(strings.Join(splitQuery, " "))
				return query, searchtypes.SubServices, awsService, false
			} else {
				return "", searchtypes.Services, awsService, false
			}
		}
	}

	return query, searchtypes.Services, nil, false
}

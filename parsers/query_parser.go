package parsers

import (
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchtypes"
)

func ParseQuery(awsServices []awsworkflow.AwsService, query string) (string, searchtypes.SearchType, *awsworkflow.AwsService) {
	// TODO add better lexing here to route populators

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
			if strings.HasPrefix(query, "$") && searchType != 0 {
				query = query[1:]
				return query, searchType, awsService
			} else {
				// prepend the home to the sub-service list so that it's still accessible
				awsServiceHome := *awsService
				awsServiceHome.Id = "home"
				awsService.SubServices = append(
					[]awsworkflow.AwsService{
						awsServiceHome,
					},
					awsService.SubServices...,
				)

				if len(awsService.SubServices) > 1 {
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
								return query, searchType, subService
							}
						}
					}
				}
				query = strings.TrimSpace(strings.Join(splitQuery, " "))
				return query, searchtypes.SubServices, awsService
			}
		}
	}

	return query, searchtypes.Services, nil
}

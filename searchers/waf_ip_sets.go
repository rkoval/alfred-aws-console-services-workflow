package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type WAFIPSetSearcher struct{}

func (s WAFIPSetSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadWafv2IPSetSummaryArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
	}
	return nil
}

func (s WAFIPSetSearcher) fetch(cfg aws.Config) ([]types.IPSetSummary, error) {
	client := wafv2.NewFromConfig(cfg)

	NextMarker := ""
	entities := []types.IPSetSummary{}
	for {
		params := &wafv2.ListIPSetsInput{
			Limit: aws.Int32(100),          // get as many as we can
			Scope: types.Scope("REGIONAL"), // TODO support CLOUDFRONT Scope somehow
		}
		if NextMarker != "" {
			params.NextMarker = &NextMarker
		}
		resp, err := client.ListIPSets(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.IPSets...)

		if resp.NextMarker != nil {
			NextMarker = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s WAFIPSetSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, entity types.IPSetSummary) {
	title := *entity.Name
	var subtitle string
	if entity.Description != nil {
		subtitle = *entity.Description
	}

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://console.aws.amazon.com/wafv2/homev2/ip-set/%s/%s?region=%s", *entity.Name, *entity.Id, config.Region)).
		Icon(awsworkflow.GetImageIcon("waf"))
}

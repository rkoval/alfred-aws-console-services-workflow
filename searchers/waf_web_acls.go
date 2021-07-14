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

type WAFWebACLSearcher struct{}

func (s WAFWebACLSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadWafv2WebACLSummaryArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
	}
	return nil
}

func (s WAFWebACLSearcher) fetch(cfg aws.Config) ([]types.WebACLSummary, error) {
	client := wafv2.NewFromConfig(cfg)

	NextMarker := ""
	entities := []types.WebACLSummary{}
	for {
		params := &wafv2.ListWebACLsInput{
			Limit: aws.Int32(100),          // get as many as we can
			Scope: types.Scope("REGIONAL"), // TODO support CLOUDFRONT Scope somehow
		}
		if NextMarker != "" {
			params.NextMarker = &NextMarker
		}
		resp, err := client.ListWebACLs(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.WebACLs...)

		if resp.NextMarker != nil {
			NextMarker = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s WAFWebACLSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, entity types.WebACLSummary) {
	title := *entity.Name
	subtitle := *entity.Description

	path := fmt.Sprintf("/wafv2/homev2/web-acl/%s/%s/overview?region=%s", *entity.Name, *entity.Id, config.Region)
	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, config.Region)).
		Icon(awsworkflow.GetImageIcon("waf"))
}

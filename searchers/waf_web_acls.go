package searchers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/wafv2"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type WAFWebACLSearcher struct{}

func (s WAFWebACLSearcher) Search(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadWafv2WebACLSummaryArrayFromCache(wf, session, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, session.Config, entity)
	}
	return nil
}

func (s WAFWebACLSearcher) fetch(session *session.Session) ([]wafv2.WebACLSummary, error) {
	client := wafv2.New(session)

	NextMarker := ""
	entities := []wafv2.WebACLSummary{}
	for {
		params := &wafv2.ListWebACLsInput{
			Limit: aws.Int64(100),         // get as many as we can
			Scope: aws.String("REGIONAL"), // TODO support CLOUDFRONT Scope somehow
		}
		if NextMarker != "" {
			params.SetNextMarker(NextMarker)
		}
		resp, err := client.ListWebACLs(params)

		if err != nil {
			return nil, err
		}

		for _, entity := range resp.WebACLs {
			entities = append(entities, *entity)
		}

		if resp.NextMarker != nil {
			NextMarker = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s WAFWebACLSearcher) addToWorkflow(wf *aw.Workflow, query string, config *aws.Config, entity wafv2.WebACLSummary) {
	title := *entity.Name
	subtitle := *entity.Description

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://console.aws.amazon.com/wafv2/homev2/web-acl/%s/%s/overview?region=%s", *entity.Name, *entity.Id, *config.Region)).
		Icon(awsworkflow.GetImageIcon("waf"))
}

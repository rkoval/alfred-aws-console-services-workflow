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

func SearchWAFIPSets(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadWafv2IPSetSummaryArrayFromCache(wf, session, cacheName, fetchIPSets, forceFetch, fullQuery)
	for _, entity := range entities {
		addIPSetToWorkflow(wf, query, session.Config, entity)
	}
	return nil
}

func fetchIPSets(session *session.Session) ([]wafv2.IPSetSummary, error) {
	client := wafv2.New(session)

	NextMarker := ""
	ipSets := []wafv2.IPSetSummary{}
	for {
		params := &wafv2.ListIPSetsInput{
			Limit: aws.Int64(100),         // get as many as we can
			Scope: aws.String("REGIONAL"), // TODO support CLOUDFRONT Scope somehow
		}
		if NextMarker != "" {
			params.SetNextMarker(NextMarker)
		}
		resp, err := client.ListIPSets(params)

		if err != nil {
			return nil, err
		}

		for _, ipSet := range resp.IPSets {
			ipSets = append(ipSets, *ipSet)
		}

		if resp.NextMarker != nil {
			NextMarker = *resp.NextMarker
		} else {
			break
		}
	}

	return ipSets, nil
}

func addIPSetToWorkflow(wf *aw.Workflow, query string, config *aws.Config, ipSet wafv2.IPSetSummary) {
	title := *ipSet.Name
	subtitle := *ipSet.Description

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://console.aws.amazon.com/wafv2/homev2/ip-set/%s/%s?region=%s", *ipSet.Name, *ipSet.Id, *config.Region)).
		Icon(awsworkflow.GetImageIcon("waf"))
}

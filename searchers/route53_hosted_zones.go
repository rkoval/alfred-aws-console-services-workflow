package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type Route53HostedZoneSearcher struct{}

func (s Route53HostedZoneSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadRoute53HostedZoneArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s Route53HostedZoneSearcher) fetch(cfg aws.Config) ([]types.HostedZone, error) {
	client := route53.NewFromConfig(cfg)

	entities := []types.HostedZone{}
	pageToken := ""
	for {
		params := &route53.ListHostedZonesInput{
			MaxItems: aws.Int32(100),
		}
		if pageToken != "" {
			params.Marker = &pageToken
		}
		resp, err := client.ListHostedZones(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.HostedZones...)

		if resp.NextMarker != nil {
			pageToken = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s Route53HostedZoneSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.HostedZone) {
	title := *entity.Name

	subtitleArray := []string{}
	if entity.Config != nil {
		config := *entity.Config

		privateString := "Public"
		if config.PrivateZone {
			privateString = "Private"
		}
		subtitleArray = util.AppendString(subtitleArray, &privateString)

		subtitleArray = util.AppendString(subtitleArray, config.Comment)
	}

	if entity.ResourceRecordSetCount != nil {
		recordSetCount := *entity.ResourceRecordSetCount
		recordSetString := fmt.Sprint(recordSetCount) + " record"
		if recordSetCount > 1 {
			recordSetString += "s"
		}
		subtitleArray = util.AppendString(subtitleArray, &recordSetString)
	}
	subtitleArray = util.AppendString(subtitleArray, entity.Id)
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/route53/v2/hostedzones#ListRecordSets/%s", *entity.Id)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.Cfg.Region)).
		Icon(awsworkflow.GetImageIcon("route53")).
		Valid(true)

	if strings.HasPrefix(searchArgs.Query, "Z") {
		item.Match(*entity.Id)
	}

}

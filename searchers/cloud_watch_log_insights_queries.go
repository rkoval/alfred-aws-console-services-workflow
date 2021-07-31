package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CloudWatchLogInsightsQuerySearcher struct{}

func (s CloudWatchLogInsightsQuerySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadCloudwatchlogsQueryDefinitionArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s CloudWatchLogInsightsQuerySearcher) fetch(cfg aws.Config) ([]types.QueryDefinition, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)

	pageToken := ""
	entities := []types.QueryDefinition{}
	for {
		params := &cloudwatchlogs.DescribeQueryDefinitionsInput{
			MaxResults: aws.Int32(1000),
		}
		if pageToken != "" {
			params.NextToken = aws.String(pageToken)
		}
		resp, err := client.DescribeQueryDefinitions(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.QueryDefinitions...)

		if resp.NextToken != nil {
			pageToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CloudWatchLogInsightsQuerySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.QueryDefinition) {
	title := *entity.Name

	subtitleArray := []string{}
	if entity.QueryString != nil {
		replaced := strings.ReplaceAll(*entity.QueryString, "\n", " ")
		subtitleArray = util.AppendString(subtitleArray, &replaced)
	}
	subtitle := strings.Join(subtitleArray, " – ")

	queryDetail := util.ConstructCloudwatchInsightsQueryDetail(entity)
	path := fmt.Sprintf("/cloudwatch/home?region=%s#logsV2:logs-insights$3FqueryDetail$3D%s", searchArgs.Cfg.Region, queryDetail)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.Cfg.Region)).
		Icon(awsworkflow.GetImageIcon("cloudwatch")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}

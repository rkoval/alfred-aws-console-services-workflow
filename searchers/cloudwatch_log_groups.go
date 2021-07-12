package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CloudWatchLogGroupSearcher struct{}

func (s CloudWatchLogGroupSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadCloudwatchlogsLogGroupArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
	}
	return nil
}

func (s CloudWatchLogGroupSearcher) fetch(cfg aws.Config) ([]types.LogGroup, error) {
	svc := cloudwatchlogs.NewFromConfig(cfg)

	NextToken := ""
	var entities []types.LogGroup
	for {
		params := &cloudwatchlogs.DescribeLogGroupsInput{
			Limit: aws.Int32(50), // get as many as we can
		}
		if NextToken != "" {
			params.NextToken = aws.String(NextToken)
		}
		resp, err := svc.DescribeLogGroups(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		for _, entity := range resp.LogGroups {
			entities = append(entities, entity)
		}

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CloudWatchLogGroupSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, entity types.LogGroup) {
	title := *entity.LogGroupName
	subtitleArray := []string{}
	if entity.StoredBytes != nil {
		subtitleArray = append(subtitleArray, fmt.Sprintf("%s stored", util.ByteFormat(*entity.StoredBytes, 2)))
	}
	if entity.RetentionInDays != nil {
		subtitleArray = append(subtitleArray, fmt.Sprintf("%d day retention", *entity.RetentionInDays))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/cloudwatch/home?region=%s#logsV2:log-groups/log-group/%s/log-events", config.Region, config.Region, url.PathEscape(*entity.LogGroupName))).
		Icon(awsworkflow.GetImageIcon("cloudwatch"))
}

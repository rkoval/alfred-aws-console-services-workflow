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
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CloudWatchLogGroupSearcher struct{}

func (s CloudWatchLogGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadCloudwatchlogsLogGroupArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
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

		entities = append(entities, resp.LogGroups...)

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CloudWatchLogGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.LogGroup) {
	title := *entity.LogGroupName
	subtitleArray := []string{}
	if entity.StoredBytes != nil {
		subtitleArray = append(subtitleArray, fmt.Sprintf("%s stored", util.ByteFormat(*entity.StoredBytes, 2)))
	}
	if entity.RetentionInDays != nil {
		subtitleArray = append(subtitleArray, fmt.Sprintf("%d day retention", *entity.RetentionInDays))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/cloudwatch/home?region=%s#logsV2:log-groups/log-group/%s/log-events", searchArgs.Cfg.Region, url.PathEscape(*entity.LogGroupName))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.Cfg.Region)).
		Icon(awsworkflow.GetImageIcon("cloudwatch"))

	searchArgs.AddMatch(item, "arn:", *entity.Arn, title)
}

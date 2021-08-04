package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SNSTopicSearcher struct{}

func (s SNSTopicSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadSnsTopicArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s SNSTopicSearcher) fetch(cfg aws.Config) ([]types.Topic, error) {
	client := sns.NewFromConfig(cfg)

	entities := []types.Topic{}
	pageToken := ""
	for {
		params := &sns.ListTopicsInput{}
		if pageToken != "" {
			params.NextToken = &pageToken
		}
		resp, err := client.ListTopics(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Topics...)

		if resp.NextToken != nil {
			pageToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s SNSTopicSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Topic) {
	subtitle := *entity.TopicArn
	title := util.GetEndOfArn(*entity.TopicArn)

	path := fmt.Sprintf("/sns/v3/home#/topic/%s", *entity.TopicArn)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("sns")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.TopicArn, title)
}

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
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SNSTopicSearcher struct{}

func (s SNSTopicSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadSnsTopicArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
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

		for _, entity := range resp.Topics {
			entities = append(entities, entity)
		}

		if resp.NextToken != nil {
			pageToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s SNSTopicSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, entity types.Topic) {
	subtitle := *entity.TopicArn
	title := util.GetEndOfArn(*entity.TopicArn)

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf(
			"https://%s.console.aws.amazon.com/sns/v3/home?region=%s#/topic/%s",
			config.Region,
			config.Region,
			*entity.TopicArn,
		)).
		Icon(awsworkflow.GetImageIcon("sns")).
		Valid(true)
}

package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SNSSubscriptionSearcher struct{}

func (s SNSSubscriptionSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadSnsSubscriptionArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
	}
	return nil
}

func (s SNSSubscriptionSearcher) fetch(cfg aws.Config) ([]types.Subscription, error) {
	client := sns.NewFromConfig(cfg)

	entities := []types.Subscription{}
	pageToken := ""
	for {
		params := &sns.ListSubscriptionsInput{}
		if pageToken != "" {
			params.NextToken = &pageToken
		}
		resp, err := client.ListSubscriptions(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		for _, entity := range resp.Subscriptions {
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

func (s SNSSubscriptionSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, entity types.Subscription) {
	topicName := util.GetEndOfArn(*entity.TopicArn)
	title := topicName

	isPending := entity.SubscriptionArn == nil || *entity.SubscriptionArn == "PendingConfirmation"
	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.Protocol)
	subtitleArray = util.AppendString(subtitleArray, entity.Endpoint)
	var subtitle string

	var url string
	if isPending {
		// subscription is still pending, so there's no permalink to it yet
		url = fmt.Sprintf(
			"https://%s.console.aws.amazon.com/sns/v3/home?region=%s#/subscriptions",
			config.Region,
			config.Region,
		)
		subtitle = "🕘 " + subtitle
	} else {
		url = fmt.Sprintf(
			"https://%s.console.aws.amazon.com/sns/v3/home?region=%s#/subscription/%s",
			config.Region,
			config.Region,
			*entity.SubscriptionArn,
		)
		subtitle = "✅ " + subtitle
		subscriptionId := util.GetEndOfArn(*entity.SubscriptionArn)
		subtitleArray = util.AppendString(subtitleArray, &subscriptionId)
	}

	subtitle += strings.Join(subtitleArray, " – ")

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(url).
		Icon(awsworkflow.GetImageIcon("sns")).
		Valid(true)
}

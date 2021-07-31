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
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SNSSubscriptionSearcher struct{}

func (s SNSSubscriptionSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadSnsSubscriptionArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
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

		entities = append(entities, resp.Subscriptions...)

		if resp.NextToken != nil {
			pageToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s SNSSubscriptionSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Subscription) {
	topicName := util.GetEndOfArn(*entity.TopicArn)
	title := topicName

	isPending := entity.SubscriptionArn == nil || *entity.SubscriptionArn == "PendingConfirmation"
	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.Protocol)
	subtitleArray = util.AppendString(subtitleArray, entity.Endpoint)
	var subtitle string

	var path string
	if isPending {
		// subscription is still pending, so there's no permalink to it yet
		path = fmt.Sprintf(
			"/sns/v3/home?region=%s#/subscriptions",
			searchArgs.Cfg.Region,
		)
		subtitle = "🕘 " + subtitle
	} else {
		path = fmt.Sprintf(
			"/sns/v3/home?region=%s#/subscription/%s",
			searchArgs.Cfg.Region,
			*entity.SubscriptionArn,
		)
		subtitle = "✅ " + subtitle
		subscriptionId := util.GetEndOfArn(*entity.SubscriptionArn)
		subtitleArray = util.AppendString(subtitleArray, &subscriptionId)
	}

	subtitle += strings.Join(subtitleArray, " – ")

	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.Cfg.Region)).
		Icon(awsworkflow.GetImageIcon("sns")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.TopicArn, title)
}

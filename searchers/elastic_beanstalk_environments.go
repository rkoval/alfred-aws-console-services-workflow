package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type ElasticBeanstalkEnvironmentSearcher struct{}

func (s ElasticBeanstalkEnvironmentSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadElasticbeanstalkEnvironmentDescriptionArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s ElasticBeanstalkEnvironmentSearcher) fetch(cfg aws.Config) ([]types.EnvironmentDescription, error) {
	svc := elasticbeanstalk.NewFromConfig(cfg)

	NextToken := ""
	environments := []types.EnvironmentDescription{}
	for {
		resp, err := svc.DescribeEnvironments(context.TODO(), &elasticbeanstalk.DescribeEnvironmentsInput{
			MaxRecords: aws.Int32(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		})
		if err != nil {
			return nil, err
		}

		for i := range resp.Environments {
			environments = append(environments, resp.Environments[i])
		}

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return environments, nil
}

func (s ElasticBeanstalkEnvironmentSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, environment types.EnvironmentDescription) {
	title := *environment.EnvironmentName
	subtitle := util.GetElasticBeanstalkHealthEmoji(environment.Health) + " " + *environment.EnvironmentId + " " + *environment.ApplicationName
	var page string
	if environment.Status == types.EnvironmentStatusTerminated {
		// "dashboard" page does not exist for terminated instances
		page = "events"
	} else {
		page = "dashboard"
	}

	path := fmt.Sprintf("/elasticbeanstalk/home?region=%s#/environment/%s?applicationName=%s&environmentId=%s", searchArgs.Cfg.Region, page, *environment.ApplicationName, *environment.EnvironmentId)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.Cfg.Region)).
		Icon(awsworkflow.GetImageIcon("elasticbeanstalk"))

	if strings.HasPrefix(searchArgs.Query, "e-") {
		item.Match(*environment.EnvironmentId)
	}
}

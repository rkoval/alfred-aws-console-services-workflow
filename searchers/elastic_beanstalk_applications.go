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

type ElasticBeanstalkApplicationSearcher struct{}

func (s ElasticBeanstalkApplicationSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s ElasticBeanstalkApplicationSearcher) fetch(cfg aws.Config) ([]types.ApplicationDescription, error) {
	client := elasticbeanstalk.NewFromConfig(cfg)

	entities := []types.ApplicationDescription{}

	params := &elasticbeanstalk.DescribeApplicationsInput{}

	resp, err := client.DescribeApplications(context.TODO(), params)

	if err != nil {
		return nil, err
	}

	entities = append(entities, resp.Applications...)

	return entities, nil
}

func (s ElasticBeanstalkApplicationSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.ApplicationDescription) {
	title := ""
	if entity.ApplicationName != nil {
		title = *entity.ApplicationName
	} else {
		title = *entity.ApplicationArn
	}

	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.Description)
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/elasticbeanstalk/home#/application/overview?applicationName=%s", *entity.ApplicationName)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("elasticbeanstalk")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.ApplicationArn, title)
}

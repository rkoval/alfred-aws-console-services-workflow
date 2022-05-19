package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EC2InstanceSearcher struct{}

func (s EC2InstanceSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EC2InstanceSearcher) fetch(cfg aws.Config) ([]types.Instance, error) {
	svc := ec2.NewFromConfig(cfg)

	NextToken := ""
	var instances []types.Instance
	for {
		params := &ec2.DescribeInstancesInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		}
		resp, err := svc.DescribeInstances(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		for _, reservation := range resp.Reservations {
			instances = append(instances, reservation.Instances...)
		}

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return instances, nil
}

func (s EC2InstanceSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Instance) {
	var title string
	subtitle := util.GetEC2InstanceStateEmoji(*entity.State)
	name := util.GetEC2TagValue(entity.Tags, "Name")
	if name != "" {
		title = name
		subtitle += " " + *entity.InstanceId
	} else {
		title = *entity.InstanceId
	}
	subtitle += " " + string(entity.InstanceType)

	path := fmt.Sprintf("/ec2/v2/home#InstanceDetails:instanceId=%s", *entity.InstanceId)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2"))

	searchArgs.AddMatch(item, "i-", *entity.InstanceId, title)
}

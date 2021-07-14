package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EC2InstanceSearcher struct{}

func (s EC2InstanceSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEc2InstanceArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
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
			for i := range reservation.Instances {
				instances = append(instances, reservation.Instances[i])
			}
		}

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return instances, nil
}

func (s EC2InstanceSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, instance types.Instance) {
	var title string
	subtitle := util.GetEC2InstanceStateEmoji(*instance.State)
	name := util.GetEC2TagValue(instance.Tags, "Name")
	if name != "" {
		title = name
		subtitle += " " + *instance.InstanceId
	} else {
		title = *instance.InstanceId
	}
	subtitle += " " + string(instance.InstanceType)

	path := fmt.Sprintf("/ec2/v2/home?region=%s#Instances:search=%s", config.Region, *instance.InstanceId)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, config.Region)).
		Icon(awsworkflow.GetImageIcon("ec2"))

	if strings.HasPrefix(query, "i-") {
		item.Match(*instance.InstanceId)
	}
}

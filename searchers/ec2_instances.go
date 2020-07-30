package searchers

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EC2InstanceSearcher struct{}

func (s EC2InstanceSearcher) Search(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	instances := caching.LoadEc2InstanceArrayFromCache(wf, session, cacheName, s.fetch, forceFetch, fullQuery)
	for _, instance := range instances {
		s.addToWorkflow(wf, query, session.Config, instance)
	}
	return nil
}

func (s EC2InstanceSearcher) fetch(session *session.Session) ([]ec2.Instance, error) {
	svc := ec2.New(session)

	NextToken := ""
	var instances []ec2.Instance
	for {
		params := &ec2.DescribeInstancesInput{
			MaxResults: aws.Int64(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		}
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			return nil, err
		}

		for _, reservation := range resp.Reservations {
			for i := range reservation.Instances {
				instances = append(instances, *reservation.Instances[i])
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

func (s EC2InstanceSearcher) addToWorkflow(wf *aw.Workflow, query string, config *aws.Config, instance ec2.Instance) {
	var title string
	subtitle := util.GetEC2InstanceStateEmoji(*instance.State)
	name := util.GetEC2TagValue(instance.Tags, "Name")
	if name != "" {
		title = name
		subtitle += " " + *instance.InstanceId
	} else {
		title = *instance.InstanceId
	}
	subtitle += " " + *instance.InstanceType

	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#Instances:search=%s", *config.Region, *config.Region, *instance.InstanceId)).
		Icon(awsworkflow.GetImageIcon("ec2"))

	if strings.HasPrefix(query, "i-") {
		item.Match(*instance.InstanceId)
	}
}

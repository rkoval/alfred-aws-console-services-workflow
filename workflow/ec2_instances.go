package workflow

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

func GetInstanceStateEmoji(instanceState string) string {
	if instanceState == ec2.InstanceStateNamePending {
		return "⚪️"
	} else if instanceState == ec2.InstanceStateNameRunning {
		return "🟢"
	} else if instanceState == ec2.InstanceStateNameShuttingDown || instanceState == ec2.InstanceStateNameStopping {
		return "🟡"
	} else if instanceState == ec2.InstanceStateNameTerminated || instanceState == ec2.InstanceStateNameStopped {
		return "🔴"
	}
	return "❔"
}

func PopulateEC2Instances(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	instances := LoadEc2InstanceArrayFromCache(wf, session, "ec2_instances", fetchEC2Instances, forceFetch, fullQuery)
	for _, instance := range instances {
		addInstanceToWorkflow(wf, query, "us-west-2" /* TODO make this read from config */, instance)
	}
	return nil
}

func fetchEC2Instances(session *session.Session) ([]ec2.Instance, error) {
	svc := ec2.New(session)

	NextToken := ""
	instances := []ec2.Instance{}
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

func addInstanceToWorkflow(wf *aw.Workflow, query, region string, instance ec2.Instance) {
	var title string
	subtitle := GetInstanceStateEmoji(*instance.State.Name)
	name := GetTagValue(instance.Tags, "Name")
	if name != "" {
		title = name
		subtitle += " " + *instance.InstanceId
	} else {
		title = *instance.InstanceId
	}
	subtitle += " " + *instance.InstanceType

	item := wf.NewItem(title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#Instances:search=%s", region, region, *instance.InstanceId)).
		Icon(core.GetImageIcon("ec2")).
		Valid(true)

	if strings.HasPrefix(query, "i-") {
		item.Match(*instance.InstanceId)
	}
}

func GetTagValue(tags []*ec2.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

package searchers

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
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

func SearchEC2Instances(wf *aw.Workflow, query string) error {
	sess, cfg := core.LoadAWSConfig()
	svc := ec2.New(sess, cfg)

	values := []*string{
		aws.String(strings.Join([]string{"*", query, "*"}, "")),
	}

	var name string
	if strings.HasPrefix(query, "i-") {
		// assume we're querying by ID here
		name = "instance-id"
	} else {
		name = "tag:Name"
	}
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(name),
				Values: values,
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		wf.NewItem(err.Error()).
			Icon(aw.IconError)
		return err
	}
	// log.Printf("%+v\n", *resp)

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
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

			wf.NewItem(title).
				Subtitle(subtitle).
				Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#Instances:search=%s", *cfg.Region, *cfg.Region, *instance.InstanceId)).
				Icon(core.GetImageIcon("ec2")).
				Valid(true)
		}
	}

	return nil
}

func GetTagValue(tags []*ec2.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

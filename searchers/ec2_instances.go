package searchers

import (
	"fmt"
	"net/http"
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

func SearchEC2Instances(wf *aw.Workflow, query string, transport http.RoundTripper) error {
	sess, cfg := core.LoadAWSConfig(transport)
	svc := ec2.New(sess, cfg)

	NextToken := ""
	for {
		params := &ec2.DescribeInstancesInput{
			MaxResults: aws.Int64(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		}
		if NextToken != "" {
			params.NextToken = aws.String(NextToken)
		}
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			wf.NewItem(err.Error()).
				Icon(aw.IconError)
			return err
		}
		// log.Println("resp", resp)

		addInstancesToWorkflow(wf, query, cfg, resp)

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return nil
}

func addInstancesToWorkflow(wf *aw.Workflow, query string, cfg *aws.Config, resp *ec2.DescribeInstancesOutput) {
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

			item := wf.NewItem(title).
				Subtitle(subtitle).
				Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#Instances:search=%s", *cfg.Region, *cfg.Region, *instance.InstanceId)).
				Icon(core.GetImageIcon("ec2")).
				Valid(true)

			if strings.HasPrefix(query, "i-") {
				item.Match(*instance.InstanceId)
			}
		}
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

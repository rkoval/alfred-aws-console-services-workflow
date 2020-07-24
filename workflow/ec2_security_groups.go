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

func PopulateEC2SecurityGroups(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	securityGroups := LoadEc2SecurityGroupArrayFromCache(wf, session, "ec2_security_groups", fetchEC2SecurityGroups, forceFetch, fullQuery)
	for _, securityGroup := range securityGroups {
		addSecurityGroupToWorkflow(wf, query, "us-west-2" /* TODO make this read from config */, securityGroup)
	}
	return nil
}

func fetchEC2SecurityGroups(session *session.Session) ([]ec2.SecurityGroup, error) {
	svc := ec2.New(session)

	NextToken := ""
	securityGroups := []ec2.SecurityGroup{}
	for {
		resp, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
			MaxResults: aws.Int64(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		})
		if err != nil {
			return nil, err
		}
		// log.Println("resp", resp)

		for i := range resp.SecurityGroups {
			securityGroups = append(securityGroups, *resp.SecurityGroups[i])
		}

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}
	return securityGroups, nil
}

func addSecurityGroupToWorkflow(wf *aw.Workflow, query, region string, securityGroup ec2.SecurityGroup) {
	var title string
	var subtitle string
	name := GetTagValue(securityGroup.Tags, "Name")
	if name != "" {
		title = name
		subtitle = *securityGroup.GroupId
	} else {
		title = *securityGroup.GroupId
	}

	if subtitle != "" {
		subtitle += " "
	}
	subtitle += *securityGroup.Description

	item := wf.NewItem(title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#SecurityGroups:group-id=%s", region, region, *securityGroup.GroupId)).
		Icon(core.GetImageIcon("ec2")).
		Valid(true)

	if strings.HasPrefix(query, "sg-") {
		item.Match(*securityGroup.GroupId)
	}
}

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

func SearchEC2SecurityGroups(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	securityGroups := caching.LoadEc2SecurityGroupArrayFromCache(wf, session, cacheName, fetchEC2SecurityGroups, forceFetch, fullQuery)
	for _, securityGroup := range securityGroups {
		addSecurityGroupToWorkflow(wf, query, session.Config, securityGroup)
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

func addSecurityGroupToWorkflow(wf *aw.Workflow, query string, config *aws.Config, securityGroup ec2.SecurityGroup) {
	var title string
	var subtitle string
	name := util.GetEC2TagValue(securityGroup.Tags, "Name")
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

	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#SecurityGroups:group-id=%s", *config.Region, *config.Region, *securityGroup.GroupId)).
		Icon(awsworkflow.GetImageIcon("ec2"))

	if strings.HasPrefix(query, "sg-") {
		item.Match(*securityGroup.GroupId)
	}
}

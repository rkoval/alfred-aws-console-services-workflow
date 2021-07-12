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

type EC2SecurityGroupSearcher struct{}

func (s EC2SecurityGroupSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEc2SecurityGroupArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
	}
	return nil
}

func (s EC2SecurityGroupSearcher) fetch(cfg aws.Config) ([]types.SecurityGroup, error) {
	svc := ec2.NewFromConfig(cfg)

	NextToken := ""
	securityGroups := []types.SecurityGroup{}
	for {
		resp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		})
		if err != nil {
			return nil, err
		}
		// log.Println("resp", resp)

		for i := range resp.SecurityGroups {
			securityGroups = append(securityGroups, resp.SecurityGroups[i])
		}

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}
	return securityGroups, nil
}

func (s EC2SecurityGroupSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, securityGroup types.SecurityGroup) {
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
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#SecurityGroups:group-id=%s", config.Region, config.Region, *securityGroup.GroupId)).
		Icon(awsworkflow.GetImageIcon("ec2"))

	if strings.HasPrefix(query, "sg-") {
		item.Match(*securityGroup.GroupId)
	}
}

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

type EC2SecurityGroupSearcher struct{}

func (s EC2SecurityGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EC2SecurityGroupSearcher) fetch(cfg aws.Config) ([]types.SecurityGroup, error) {
	svc := ec2.NewFromConfig(cfg)

	NextToken := ""
	entities := []types.SecurityGroup{}
	for {
		resp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		})
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.SecurityGroups...)

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}
	return entities, nil
}

func (s EC2SecurityGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.SecurityGroup) {
	var title string
	var subtitle string
	name := util.GetEC2TagValue(entity.Tags, "Name")
	if name != "" {
		title = name
		subtitle = *entity.GroupId
	} else {
		title = *entity.GroupId
	}

	if subtitle != "" {
		subtitle += " "
	}
	subtitle += *entity.Description

	path := fmt.Sprintf("/ec2/v2/home#SecurityGroups:group-id=%s", *entity.GroupId)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2"))

	searchArgs.AddMatch(item, "sg-", *entity.GroupId, title)
}

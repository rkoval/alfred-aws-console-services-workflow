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

func SearchEC2SecurityGroups(wf *aw.Workflow, query string, transport http.RoundTripper) error {
	sess, cfg := core.LoadAWSConfig(transport)
	svc := ec2.New(sess, cfg)

	values := []*string{
		aws.String(strings.Join([]string{"*", query, "*"}, "")),
	}

	var name string
	if strings.HasPrefix(query, "sg-") {
		// assume we're querying by ID here
		name = "group-id"
	} else {
		name = "tag:Name"
	}
	params := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(name),
				Values: values,
			},
		},
	}

	resp, err := svc.DescribeSecurityGroups(params)
	if err != nil {
		wf.NewItem(err.Error()).
			Icon(aw.IconError)
		return err
	}
	// log.Printf("%+v\n", *resp)

	for _, securityGroup := range resp.SecurityGroups {
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

		wf.NewItem(title).
			Subtitle(subtitle).
			Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#SecurityGroups:group-id=%s", *cfg.Region, *cfg.Region, *securityGroup.GroupId)).
			Icon(core.GetImageIcon("ec2")).
			Valid(true)
	}

	return nil
}

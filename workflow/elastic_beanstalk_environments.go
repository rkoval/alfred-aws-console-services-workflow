package workflow

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

func getHealthEmoji(environmentHealth string) string {
	if environmentHealth == elasticbeanstalk.EnvironmentHealthGreen {
		return "🟢"
	} else if environmentHealth == elasticbeanstalk.EnvironmentHealthYellow {
		return "🟡"
	} else if environmentHealth == elasticbeanstalk.EnvironmentHealthRed {
		return "🔴"
	} else if environmentHealth == elasticbeanstalk.EnvironmentHealthGrey {
		return "⚪️"
	}

	return "❔"
}

func PopulateElasticBeanstalkEnvironments(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	instances := LoadElasticbeanstalkEnvironmentDescriptionArrayFromCache(wf, session, "ec2_instances", fetchElasticBeanstalkEnvironments, forceFetch, fullQuery)
	for _, instance := range instances {
		addEnvironmentToWorkflow(wf, query, "us-west-2" /* TODO make this read from config */, instance)
	}
	return nil
}

func fetchElasticBeanstalkEnvironments(session *session.Session) ([]elasticbeanstalk.EnvironmentDescription, error) {
	svc := elasticbeanstalk.New(session)

	NextToken := ""
	environments := []elasticbeanstalk.EnvironmentDescription{}
	for {
		resp, err := svc.DescribeEnvironments(&elasticbeanstalk.DescribeEnvironmentsInput{
			MaxRecords: aws.Int64(1000), // get as many as we can
			NextToken:  aws.String(NextToken),
		})
		if err != nil {
			return nil, err
		}

		for i := range resp.Environments {
			environments = append(environments, *resp.Environments[i])
		}

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return environments, nil
}

func addEnvironmentToWorkflow(wf *aw.Workflow, query, region string, environment elasticbeanstalk.EnvironmentDescription) {
	title := *environment.EnvironmentName
	subtitle := getHealthEmoji(*environment.Health) + " " + *environment.EnvironmentId + " " + *environment.ApplicationName
	var page string
	if *environment.Status == elasticbeanstalk.EnvironmentStatusTerminated {
		// "dashboard" page does not exist for terminated instances
		page = "events"
	} else {
		page = "dashboard"
	}
	item := wf.NewItem(title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf(
			"https://%s.console.aws.amazon.com/elasticbeanstalk/home?region=%s#/environment/%s?applicationName=%s&environmentId=%s",
			region,
			region,
			page,
			*environment.ApplicationName,
			*environment.EnvironmentId,
		)).
		Icon(core.GetImageIcon("elasticbeanstalk")).
		Valid(true)

	if strings.HasPrefix(query, "e-") {
		item.Match(*environment.EnvironmentId)
	}
}

package searchers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

func GetHealthEmoji(environmentHealth string) string {
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

func SearchElasticBeanstalkEnvironments(wf *aw.Workflow, query string, transport http.RoundTripper) error {
	sess, cfg := core.LoadAWSConfig(transport)
	svc := elasticbeanstalk.New(sess, cfg)
	params := &elasticbeanstalk.DescribeEnvironmentsInput{}
	resp, err := svc.DescribeEnvironments(params)
	if err != nil {
		wf.NewItem(err.Error()).
			Icon(aw.IconError)
		return err
	}
	// log.Println("resp", resp)

	for _, environment := range resp.Environments {
		title := *environment.EnvironmentName
		subtitle := GetHealthEmoji(*environment.Health) + " " + *environment.EnvironmentId + " " + *environment.ApplicationName
		var page string
		if *environment.Status == elasticbeanstalk.EnvironmentStatusTerminated {
			// "dashboard" page does not exist for terminated instances
			page = "events"
		} else {
			page = "dashboard"
		}
		wf.NewItem(title).
			Subtitle(subtitle).
			Arg(fmt.Sprintf(
				"https://%s.console.aws.amazon.com/elasticbeanstalk/home?region=%s#/environment/%s?applicationName=%s&environmentId=%s",
				*cfg.Region,
				*cfg.Region,
				page,
				*environment.ApplicationName,
				*environment.EnvironmentId,
			)).
			Icon(core.GetImageIcon("elasticbeanstalk")).
			Valid(true)
	}

	return nil
}

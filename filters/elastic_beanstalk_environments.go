package filters

import (
	"fmt"

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

func SearchElasticBeanstalkEnvironments(wf *aw.Workflow, query string) error {
	sess, cfg := core.LoadAWSConfig()
	svc := elasticbeanstalk.New(sess, cfg)
	params := &elasticbeanstalk.DescribeEnvironmentsInput{}
	resp, err := svc.DescribeEnvironments(params)
	if err != nil {
		wf.NewItem(err.Error()).
			Icon(aw.IconError)
		return err
	}
	// log.Printf("%+v\n", *resp)

	for _, environment := range resp.Environments {
		title := *environment.EnvironmentName
		subtitle := GetHealthEmoji(*environment.Health) + " " + *environment.EnvironmentId + " " + *environment.ApplicationName
		wf.NewItem(title).
			Subtitle(subtitle).
			Arg(fmt.Sprintf(
				"https://%s.console.aws.amazon.com/elasticbeanstalk/home?region=%s#/environment/dashboard?applicationName=%s&environmentId=%s",
				*cfg.Region,
				*cfg.Region,
				*environment.ApplicationName,
				*environment.EnvironmentId,
			)).
			Icon(core.GetImageIcon("elasticbeanstalk")).
			Valid(true)
	}

	return nil
}

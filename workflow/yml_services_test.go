package workflow

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
)

var awsServices []awsworkflow.AwsService = []awsworkflow.AwsService{
	{
		Id:          "service1",
		Name:        "Service 1",
		ShortName:   "S1",
		Description: "Description of the first service",
		Url:         "/pizza",
	},
	{
		Id:          "service2",
		Name:        "Service 2",
		Description: "Description of the second service",
		Url:         "https://ryankoval.com",
		SubServices: []awsworkflow.AwsService{
			{
				Id:          "sub-service1",
				Name:        "Sub-service 1",
				Description: "Description of the first sub-service",
				Url:         "/bookmarks",
			},
		},
	},
	{
		Id:          "whoa",
		Name:        "Whoa",
		ShortName:   "W",
		Description: "Whoa!!!!",
		Url:         "/github",
	},
}

func TestSearchServices(t *testing.T) {
	wf := aw.New()

	cfg := awsworkflow.InitAWS(nil)
	SearchServices(wf, awsServices, cfg)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

func TestSearchSubServices(t *testing.T) {
	wf := aw.New()

	cfg := awsworkflow.InitAWS(nil)

	SearchSubServices(wf, awsServices[1], cfg)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

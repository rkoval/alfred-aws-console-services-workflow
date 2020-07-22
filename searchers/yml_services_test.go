package searchers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

var awsServices []core.AwsService = []core.AwsService{
	{
		Id:          "service1",
		Name:        "Service 1",
		ShortName:   "S1",
		Description: "Description of the first service",
		Url:         "https://ryankoval.pizza",
	},
	{
		Id:          "service2",
		Name:        "Service 2",
		Description: "Description of the second service",
		Url:         "https://ryankoval.com",
		Sections: []core.AwsServiceSection{
			{
				Id:          "section1",
				Name:        "Section 1",
				Description: "Description of the first section",
				Url:         "https://bookmarks.ryankoval.com",
			},
		},
	},
	{
		Id:          "whoa",
		Name:        "Whoa",
		ShortName:   "W",
		Description: "Whoa!!!!",
		Url:         "https://github.ryankoval.com",
	},
}

func TestSearchServices(t *testing.T) {
	wf := aw.New()

	SearchServices(wf, awsServices)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

func TestSearchServiceSections(t *testing.T) {
	wf := aw.New()

	SearchServiceSections(wf, awsServices[1])

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

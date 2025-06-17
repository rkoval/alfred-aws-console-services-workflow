package searchers

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearcher(t *testing.T, searcher Searcher, fixtureFilename string) {
	wf := aw.New()

	r := tests.NewAWSRecorderSession(fixtureFilename)
	defer tests.PanicOnError(r.Stop)

	cfg := awsworkflow.InitAWS(r, nil, nil, nil)
	err := searcher.Search(
		wf,
		searchutil.SearchArgs{
			Cfg:        cfg,
			ForceFetch: true,
			GetRegionFunc: func(cfg aws.Config) string {
				return cfg.Region
			},
		},
	)
	if err != nil {
		t.Errorf("got error from search: %v", err)
	}

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

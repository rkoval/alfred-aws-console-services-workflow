package searchers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearcher(t *testing.T, searcher Searcher, fixtureFilename string) {
	wf := aw.New()

	cfg, r := tests.NewAWSRecorderSession(fixtureFilename)
	defer tests.PanicOnError(r.Stop)

	err := searcher.Search(
		wf,
		searchutil.SearchArgs{
			Cfg:        cfg,
			ForceFetch: true,
		},
	)
	if err != nil {
		t.Errorf("got error from search: %v", err)
	}

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCloudWatchLogInsightsQuerySearcher(t *testing.T) {
	TestSearcher(t, CloudWatchLogInsightsQuerySearcher{}, util.GetCurrentFilename())
}

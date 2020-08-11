package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCloudWatchLogGroupSearcher(t *testing.T) {
	TestSearcher(t, CloudWatchLogGroupSearcher{}, util.GetCurrentFilename())
}

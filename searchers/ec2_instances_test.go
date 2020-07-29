package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSearchEC2Instances(t *testing.T) {
	TestSearcher(t, SearchEC2Instances, util.GetCurrentFilename())
}

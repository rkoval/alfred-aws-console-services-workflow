package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSearchEC2SecurityGroups(t *testing.T) {
	TestSearcher(t, SearchEC2SecurityGroups, util.GetCurrentFilename())
}

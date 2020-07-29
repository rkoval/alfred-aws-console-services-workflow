package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSearchWAFIPSets(t *testing.T) {
	TestSearcher(t, SearchWAFIPSets, util.GetCurrentFilename())
}

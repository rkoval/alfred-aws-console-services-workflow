package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSearchWAFWebACLs(t *testing.T) {
	TestSearcher(t, SearchWAFWebACLs, util.GetCurrentFilename())
}

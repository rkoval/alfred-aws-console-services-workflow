package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestRDSDatabaseSearcher(t *testing.T) {
	TestSearcher(t, RDSDatabaseSearcher{}, util.GetCurrentFilename())
}

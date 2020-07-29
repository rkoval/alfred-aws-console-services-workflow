package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestWAFIPSetSearcher(t *testing.T) {
	TestSearcher(t, WAFIPSetSearcher{}, util.GetCurrentFilename())
}

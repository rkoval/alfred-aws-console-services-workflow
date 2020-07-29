package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestWAFWebACLSearcher(t *testing.T) {
	TestSearcher(t, WAFWebACLSearcher{}, util.GetCurrentFilename())
}

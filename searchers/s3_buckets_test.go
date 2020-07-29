package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSearchS3Buckets(t *testing.T) {
	TestSearcher(t, SearchS3Buckets, util.GetCurrentFilename())
}

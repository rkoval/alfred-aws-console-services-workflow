package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestS3BucketSearcher(t *testing.T) {
	TestSearcher(t, S3BucketSearcher{}, util.GetCurrentFilename())
}

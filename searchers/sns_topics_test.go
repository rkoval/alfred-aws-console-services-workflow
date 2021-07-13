package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSNSTopicSearcher(t *testing.T) {
	TestSearcher(t, SNSTopicSearcher{}, util.GetCurrentFilename())
}

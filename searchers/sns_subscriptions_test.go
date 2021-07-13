package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSNSSubscriptionSearcher(t *testing.T) {
	TestSearcher(t, SNSSubscriptionSearcher{}, util.GetCurrentFilename())
}

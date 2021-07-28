package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestRoute53HostedZoneSearcher(t *testing.T) {
	TestSearcher(t, Route53HostedZoneSearcher{}, util.GetCurrentFilename())
}

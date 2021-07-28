package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEC2LoadBalancerSearcher(t *testing.T) {
	TestSearcher(t, EC2LoadBalancerSearcher{}, util.GetCurrentFilename())
}

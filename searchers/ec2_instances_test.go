package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEC2InstanceSearcher(t *testing.T) {
	TestSearcher(t, EC2InstanceSearcher{}, util.GetCurrentFilename())
}

package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEC2SecurityGroupSearcher(t *testing.T) {
	TestSearcher(t, EC2SecurityGroupSearcher{}, util.GetCurrentFilename())
}

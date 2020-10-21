package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCloudFormationStackSearcher(t *testing.T) {
	TestSearcher(t, CloudFormationStackSearcher{}, util.GetCurrentFilename())
}

package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestECSClusterSearcher(t *testing.T) {
	TestSearcher(t, ECSClusterSearcher{}, util.GetCurrentFilename())
}

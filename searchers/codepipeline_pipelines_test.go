package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCodePipelinePipelinesSearcher(t *testing.T) {
	TestSearcher(t, CodePipelinePipelinesSearcher{}, util.GetCurrentFilename())
}

package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestLambdaFunctionSearcher(t *testing.T) {
	TestSearcher(t, LambdaFunctionSearcher{}, util.GetCurrentFilename())
}

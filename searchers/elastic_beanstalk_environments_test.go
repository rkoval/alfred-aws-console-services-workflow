package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSearchElasticBeanstalkEnvironments(t *testing.T) {
	TestSearcher(t, SearchElasticBeanstalkEnvironments, util.GetCurrentFilename())
}

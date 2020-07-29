package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestElasticBeanstalkEnvironmentSearcher(t *testing.T) {
	TestSearcher(t, ElasticBeanstalkEnvironmentSearcher{}, util.GetCurrentFilename())
}

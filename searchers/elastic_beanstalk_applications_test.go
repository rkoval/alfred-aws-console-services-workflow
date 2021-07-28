package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestElasticBeanstalkApplicationSearcher(t *testing.T) {
	TestSearcher(t, ElasticBeanstalkApplicationSearcher{}, util.GetCurrentFilename())
}

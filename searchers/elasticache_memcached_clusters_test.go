package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestElasticacheMemcachedClusterSearcher(t *testing.T) {
	TestSearcher(t, ElasticacheMemcachedClusterSearcher{}, util.GetCurrentFilename())
}

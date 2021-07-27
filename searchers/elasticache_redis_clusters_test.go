package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestElasticacheRedisClusterSearcher(t *testing.T) {
	TestSearcher(t, ElasticacheRedisClusterSearcher{}, util.GetCurrentFilename())
}

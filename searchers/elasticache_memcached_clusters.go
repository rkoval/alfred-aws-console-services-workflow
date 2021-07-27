package searchers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/elasticacheutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type ElasticacheMemcachedClusterSearcher struct{}

func (s ElasticacheMemcachedClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadElasticacheCacheClusterArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s ElasticacheMemcachedClusterSearcher) fetch(cfg aws.Config) ([]types.CacheCluster, error) {
	return elasticacheutil.Fetch(cfg)
}

func (s ElasticacheMemcachedClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.CacheCluster) {
	elasticacheutil.AddCacheClusterToWorkflow("memcached", wf, searchArgs, entity)
}

package elasticacheutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Fetch(cfg aws.Config) ([]types.CacheCluster, error) {
	client := elasticache.NewFromConfig(cfg)

	entities := []types.CacheCluster{}
	pageToken := ""
	for {
		params := &elasticache.DescribeCacheClustersInput{
			MaxRecords: aws.Int32(100),
		}
		if pageToken != "" {
			params.Marker = &pageToken
		}
		resp, err := client.DescribeCacheClusters(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.CacheClusters...)

		if resp.Marker != nil {
			pageToken = *resp.Marker
		} else {
			break
		}
	}

	return entities, nil
}

func AddCacheClusterToWorkflow(engineName string, wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.CacheCluster) {
	if entity.Engine == nil || *entity.Engine != engineName {
		return
	}
	var title string
	if entity.CacheClusterId != nil {
		title = *entity.CacheClusterId
	} else {
		title = *entity.ARN
	}

	subtitle := util.GetElasticacheCacheClusterSubtitle(entity)

	location := engineName + "-detail"
	if entity.ReplicationGroupId != nil {
		location = engineName + "-group-detail"
	}
	path := fmt.Sprintf("/elasticache/home#%s:id=%s", location, *entity.CacheClusterId)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("elasticache")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.ARN, title)
}

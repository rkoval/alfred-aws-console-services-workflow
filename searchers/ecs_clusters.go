package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type ECSClusterSearcher struct{}

func (s ECSClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s ECSClusterSearcher) fetch(cfg aws.Config) ([]types.Cluster, error) {
	svc := ecs.NewFromConfig(cfg)

	var clusterARNs []string
	paginator := ecs.NewListClustersPaginator(svc, &ecs.ListClustersInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		clusterARNs = append(clusterARNs, output.ClusterArns...)
	}

	if len(clusterARNs) == 0 {
		return []types.Cluster{}, nil
	}

	var clusters []types.Cluster
	// DescribeClusters has a limit of 100 clusters per call
	chunkSize := 100
	for i := 0; i < len(clusterARNs); i += chunkSize {
		end := i + chunkSize
		if end > len(clusterARNs) {
			end = len(clusterARNs)
		}
		chunk := clusterARNs[i:end]

		descResp, err := svc.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
			Clusters: chunk,
			Include: []types.ClusterField{
				types.ClusterFieldTags, // Include tags
			},
		})
		if err != nil {
			// Log or handle error, but continue processing other chunks if possible
			fmt.Printf("Error describing clusters chunk %d: %v\n", i/chunkSize+1, err) // Consider better logging
			continue
		}
		clusters = append(clusters, descResp.Clusters...)
	}

	return clusters, nil
}

func (s ECSClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Cluster) {
	title := *entity.ClusterName
	subtitle := *entity.ClusterArn

	// Extract cluster name from ARN for the path
	arnParts := strings.Split(*entity.ClusterArn, "/")
	clusterName := arnParts[len(arnParts)-1]

	path := fmt.Sprintf("/ecs/v2/clusters/%s/services?region=%s", clusterName, searchArgs.GetRegion()) // Verify correct path format
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ecs")) // Ensure "ecs" icon exists

	searchArgs.AddMatch(item, "arn:", *entity.ClusterArn, title)
}

package searchers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CodePipelinePipelinesSearcher struct{}

func (s CodePipelinePipelinesSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s CodePipelinePipelinesSearcher) fetch(cfg aws.Config) ([]types.PipelineSummary, error) {
	svc := codepipeline.NewFromConfig(cfg)

	NextToken := ""
	var entities []types.PipelineSummary
	for {
		params := &codepipeline.ListPipelinesInput{
			MaxResults: aws.Int32(100),
		}
		if NextToken != "" {
			params.NextToken = aws.String(NextToken)
		}
		resp, err := svc.ListPipelines(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Pipelines...)

		if resp.NextToken != nil {
			NextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CodePipelinePipelinesSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.PipelineSummary) {
	title := *entity.Name
	subtitleArray := []string{}
	if entity.Version != nil {
		subtitleArray = append(subtitleArray, fmt.Sprintf("Version %d", *entity.Version))
	}
	if entity.Created != nil {
		subtitleArray = append(subtitleArray, fmt.Sprintf("Created %s", entity.Created.Format(time.UnixDate)))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/codesuite/codepipeline/pipelines/%s/view", title)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("codepipeline"))

	searchArgs.AddMatch(item, "", "", title)
}

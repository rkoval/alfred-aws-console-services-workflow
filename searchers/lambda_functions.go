package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type LambdaFunctionSearcher struct{}

func (s LambdaFunctionSearcher) Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadLambdaFunctionConfigurationArrayFromCache(wf, cfg, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, cfg, entity)
	}
	return nil
}

func (s LambdaFunctionSearcher) fetch(cfg aws.Config) ([]types.FunctionConfiguration, error) {
	svc := lambda.NewFromConfig(cfg)

	NextMarker := ""
	var entities []types.FunctionConfiguration
	for {
		params := &lambda.ListFunctionsInput{
			MaxItems: aws.Int32(200), // get as many as we can
		}
		if NextMarker != "" {
			params.Marker = aws.String(NextMarker)
		}
		resp, err := svc.ListFunctions(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Functions...)

		if resp.NextMarker != nil {
			NextMarker = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s LambdaFunctionSearcher) addToWorkflow(wf *aw.Workflow, query string, config aws.Config, entity types.FunctionConfiguration) {
	title := *entity.FunctionName
	subtitleArray := []string{}
	if entity.Description != nil && *entity.Description != "" {
		subtitleArray = append(subtitleArray, *entity.Description)
	}
	if entity.Runtime != "" {
		subtitleArray = append(subtitleArray, string(entity.Runtime))
	}
	if entity.CodeSize != 0 {
		subtitleArray = append(subtitleArray, util.ByteFormat(entity.CodeSize, 2))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/lambda/home?region=%s#/functions/%s?tab=configuration", config.Region, config.Region, url.PathEscape(*entity.FunctionName))).
		Icon(awsworkflow.GetImageIcon("lambda"))
}

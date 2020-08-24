package searchers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type LambdaFunctionSearcher struct{}

func (s LambdaFunctionSearcher) Search(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadLambdaFunctionConfigurationArrayFromCache(wf, session, cacheName, s.fetch, forceFetch, fullQuery)
	for _, entity := range entities {
		s.addToWorkflow(wf, query, session.Config, entity)
	}
	return nil
}

func (s LambdaFunctionSearcher) fetch(session *session.Session) ([]lambda.FunctionConfiguration, error) {
	svc := lambda.New(session)

	NextMarker := ""
	var entities []lambda.FunctionConfiguration
	for {
		params := &lambda.ListFunctionsInput{
			MaxItems: aws.Int64(200), // get as many as we can
		}
		if NextMarker != "" {
			params.Marker = aws.String(NextMarker)
		}
		resp, err := svc.ListFunctions(params)
		if err != nil {
			return nil, err
		}

		for _, entity := range resp.Functions {
			entities = append(entities, *entity)
		}

		if resp.NextMarker != nil {
			NextMarker = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s LambdaFunctionSearcher) addToWorkflow(wf *aw.Workflow, query string, config *aws.Config, entity lambda.FunctionConfiguration) {
	title := *entity.FunctionName
	subtitleArray := []string{}
	if entity.Description != nil && *entity.Description != "" {
		subtitleArray = append(subtitleArray, *entity.Description)
	}
	if entity.Runtime != nil {
		subtitleArray = append(subtitleArray, *entity.Runtime)
	}
	if entity.CodeSize != nil {
		subtitleArray = append(subtitleArray, util.ByteFormat(*entity.CodeSize, 2))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf("https://%s.console.aws.amazon.com/lambda/home?region=%s#/functions/%s?tab=configuration", *config.Region, *config.Region, url.PathEscape(*entity.FunctionName))).
		Icon(awsworkflow.GetImageIcon("lambda"))
}

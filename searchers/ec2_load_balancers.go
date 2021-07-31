package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EC2LoadBalancerSearcher struct{}

func (s EC2LoadBalancerSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadElasticloadbalancingv2LoadBalancerArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EC2LoadBalancerSearcher) fetch(cfg aws.Config) ([]types.LoadBalancer, error) {
	client := elasticloadbalancingv2.NewFromConfig(cfg)

	entities := []types.LoadBalancer{}
	pageToken := ""
	for {
		params := &elasticloadbalancingv2.DescribeLoadBalancersInput{
			PageSize: aws.Int32(400),
		}
		if pageToken != "" {
			params.Marker = &pageToken
		}
		resp, err := client.DescribeLoadBalancers(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.LoadBalancers...)

		if resp.NextMarker != nil {
			pageToken = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s EC2LoadBalancerSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.LoadBalancer) {
	title := ""
	if entity.LoadBalancerName != nil {
		title = *entity.LoadBalancerName
	} else {
		title = *entity.LoadBalancerArn
	}

	subtitleArray := []string{}
	typeString := string(entity.Type)
	subtitleArray = util.AppendString(subtitleArray, &typeString)
	subtitleArray = util.AppendString(subtitleArray, entity.DNSName)
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/ec2/home?region=%s#LoadBalancers:search=%s;sort=loadBalancerName", searchArgs.Cfg.Region, *entity.LoadBalancerArn)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.Cfg.Region)).
		Icon(awsworkflow.GetImageIcon("ec2")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.LoadBalancerArn, title)
}

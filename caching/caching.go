package caching

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudformation "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	codepipeline "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecs "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	elasticache "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elasticbeanstalk "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	elasticloadbalancingv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	lambda "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	rds "github.com/aws/aws-sdk-go-v2/service/rds/types"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53/types"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3/types"
	sns "github.com/aws/aws-sdk-go-v2/service/sns/types"
	wafv2 "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/aws/smithy-go"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type Entity interface {
	cloudwatchlogs.LogGroup | ec2.Instance | s3.Bucket | ec2.SecurityGroup | elasticbeanstalk.EnvironmentDescription | wafv2.IPSetSummary | wafv2.WebACLSummary | lambda.FunctionConfiguration | cloudformation.Stack | rds.DBInstance | sns.Topic | sns.Subscription | elasticache.CacheCluster | elasticloadbalancingv2.LoadBalancer | elasticbeanstalk.ApplicationDescription | route53.HostedZone | cloudwatchlogs.QueryDefinition | codepipeline.PipelineSummary | ecs.Cluster
}

func LoadEntityArrayFromCache[K Entity](wf *aw.Workflow, searchArgs searchutil.SearchArgs, cacheName string, fetcher func(aws.Config) ([]K, error)) []K {
	// TODO optimization: not all services have sa region associated with them, so cache can be reused across regions (e.g., s3 buckets are global)
	cacheName += "_" + searchArgs.Cfg.Region + "_" + searchArgs.Profile

	results := []K{}
	lastFetchErrPath := wf.CacheDir() + "/last-fetch-err.txt"
	if searchArgs.ForceFetch {
		log.Printf("fetching from aws ...")
		results, err := fetcher(searchArgs.Cfg)

		if err != nil {
			log.Printf("fetch error occurred. writing to %s ...", lastFetchErrPath)
			var errString string
			var missingRegionError *aws.MissingRegionError
			if errors.As(err, &missingRegionError) {
				errString = "MissingRegion"
			} else {
				var apiErr smithy.APIError
				if errors.As(err, &apiErr) {
					errCode := apiErr.ErrorCode()
					if errCode == "AccessDeniedException" {
						errString = "You do not have access to fetch these. Check your IAM permissions"
					} else {
						errString = errCode

						message := apiErr.ErrorMessage()
						if message != "" {
							errString += ": " + message
						}
					}
				}
			}
			if errString == "" {
				errString = err.Error()
			}
			if strings.Contains(errString, "failed to retrieve credentials") {
				// workaround hack; aws-sdk-go-v2 will automatically attempt to get credentials from the metadata service URL if file not specified,
				// but that's bad given that we will never be run in AWS. as a result, just populate an error string that informs users better
				errString = "NoCredentialProviders"
			}
			_ = os.WriteFile(lastFetchErrPath, []byte(errString), 0600)
			panic(err)
		} else {
			os.Remove(lastFetchErrPath)
		}
		log.Printf("fetched %d results from aws", len(results))

		log.Printf("storing %d results with cache key `%s` to %s ...", len(results), cacheName, wf.CacheDir())
		if err := wf.Cache.StoreJSON(cacheName, results); err != nil {
			panic(err)
		}
		return results
	}

	err := handleExpiredCache(wf, cacheName, lastFetchErrPath, searchArgs)
	if err != nil {
		return []K{}
	}

	if wf.Cache.Exists(cacheName) {
		log.Printf("using cache with key `%s` in %s ...", cacheName, wf.CacheDir())
		if err := wf.Cache.LoadJSON(cacheName, &results); err != nil {
			panic(err)
		}
	} else {
		log.Printf("cache with key `%s` did not exist in %s ...", cacheName, wf.CacheDir())
		wf.NewItem("Fetching ...").
			Icon(aw.IconInfo)
	}

	return results
}

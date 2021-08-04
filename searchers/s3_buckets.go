package searchers

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type S3BucketSearcher struct{}

func (s S3BucketSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	es := caching.LoadS3BucketArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range es {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s S3BucketSearcher) fetch(cfg aws.Config) ([]types.Bucket, error) {
	svc := s3.NewFromConfig(cfg)

	resp, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	buckets := []types.Bucket{}
	for i := range resp.Buckets {
		buckets = append(buckets, resp.Buckets[i])
	}
	return buckets, nil
}

func (s S3BucketSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Bucket) {
	title := *entity.Name
	subtitle := "Created " + entity.CreationDate.Format(time.UnixDate)

	// must manually append region here because wafv2 is technically a global region, but entities within it are region-specific
	path := fmt.Sprintf("/s3/buckets/%s/?region=%s&tab=objects", *entity.Name, searchArgs.Cfg.Region)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("s3")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}

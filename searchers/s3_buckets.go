package searchers

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type S3BucketSearcher struct{}

func (s S3BucketSearcher) Search(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	es := caching.LoadS3BucketArrayFromCache(wf, session, cacheName, s.fetch, forceFetch, fullQuery)
	for _, e := range es {
		s.addToWorkflow(wf, query, session.Config, e)
	}
	return nil
}

func (s S3BucketSearcher) fetch(session *session.Session) ([]s3.Bucket, error) {
	svc := s3.New(session)

	resp, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	// log.Println("resp", resp)

	buckets := []s3.Bucket{}
	for i := range resp.Buckets {
		buckets = append(buckets, *resp.Buckets[i])
	}
	return buckets, nil
}

func (s S3BucketSearcher) addToWorkflow(wf *aw.Workflow, query string, config *aws.Config, bucket s3.Bucket) {
	title := *bucket.Name
	subtitle := "Created " + bucket.CreationDate.Format(time.UnixDate)

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf(
			"https://console.aws.amazon.com/s3/buckets/%s/?region=%s&tab=overview",
			*bucket.Name,
			*config.Region,
		)).
		Icon(awsworkflow.GetImageIcon("s3")).
		Valid(true)
}

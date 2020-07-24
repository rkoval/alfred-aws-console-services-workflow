package searchers

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
)

func SearchS3Buckets(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	es := caching.LoadS3BucketArrayFromCache(wf, session, "s3_buckets", fetchS3Buckets, forceFetch, fullQuery)
	for _, e := range es {
		addS3BucketToWorkflow(wf, query, "us-west-2" /* TODO make this read from config */, e)
	}
	return nil
}

func fetchS3Buckets(session *session.Session) ([]s3.Bucket, error) {
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

func addS3BucketToWorkflow(wf *aw.Workflow, query, region string, bucket s3.Bucket) {
	title := *bucket.Name
	subtitle := "Created " + bucket.CreationDate.Format(time.UnixDate)

	wf.NewItem(title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf(
			"https://console.aws.amazon.com/s3/buckets/%s/?region=%s&tab=overview",
			*bucket.Name,
			region,
		)).
		Icon(awsworkflow.GetImageIcon("s3")).
		Valid(true)
}

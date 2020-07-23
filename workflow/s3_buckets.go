package workflow

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

func SearchS3Buckets(wf *aw.Workflow, query string, transport http.RoundTripper) error {
	sess, cfg := core.LoadAWSConfig(transport)
	svc := s3.New(sess, cfg)

	resp, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		wf.NewItem(err.Error()).
			Icon(aw.IconError)
		return err
	}
	// log.Println("resp", resp)

	for _, bucket := range resp.Buckets {
		title := *bucket.Name
		subtitle := "Created " + bucket.CreationDate.Format(time.UnixDate)

		wf.NewItem(title).
			Subtitle(subtitle).
			Arg(fmt.Sprintf(
				"https://console.aws.amazon.com/s3/buckets/%s/?region=%s&tab=overview",
				*bucket.Name,
				*cfg.Region,
			)).
			Icon(core.GetImageIcon("s3")).
			Valid(true)
	}

	return nil
}

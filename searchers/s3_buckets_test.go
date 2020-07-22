package searchers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearchS3Buckets(t *testing.T) {
	wf := aw.New()

	r := tests.NewAWSRecorder("s3_buckets_test")
	defer r.Stop()
	SearchS3Buckets(wf, "", r)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

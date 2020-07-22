package searchers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearchEC2Instances(t *testing.T) {
	wf := aw.New()

	r := tests.NewAWSRecorder("ec2_instances_fixture")
	defer r.Stop()
	SearchEC2Instances(wf, "", r)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

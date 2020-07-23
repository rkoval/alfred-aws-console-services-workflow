package workflow

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestPopulateEC2Instances(t *testing.T) {
	wf := aw.New()

	r := tests.NewAWSRecorder("ec2_instances_test")
	defer tests.PanicOnError(r.Stop)
	err := PopulateEC2Instances(wf, "", r, true, "")
	if err != nil {
		t.Errorf("got error from search: %v", err)
	}

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

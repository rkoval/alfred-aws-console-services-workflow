package workflow

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearchEC2Instances(t *testing.T) {
	wf := aw.New()

	session, r := tests.NewAWSRecorderSession("ec2_instances_test")
	defer tests.PanicOnError(r.Stop)
	err := SearchEC2Instances(wf, "", session, true, "")
	if err != nil {
		t.Errorf("got error from search: %v", err)
	}

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

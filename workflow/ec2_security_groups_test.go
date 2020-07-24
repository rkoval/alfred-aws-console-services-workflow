package workflow

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestPopulateEC2SecurityGroups(t *testing.T) {
	wf := aw.New()

	r := tests.NewAWSRecorder("ec2_security_groups_test")
	defer tests.PanicOnError(r.Stop)
	err := PopulateEC2SecurityGroups(wf, "", r, true, "")
	if err != nil {
		t.Errorf("got error from search: %v", err)
	}

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

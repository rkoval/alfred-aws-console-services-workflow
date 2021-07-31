package parsers

import (
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
)

type Query struct {
	RawQuery              string
	Service               *awsworkflow.AwsService
	SubService            *awsworkflow.AwsService
	HasTrailingWhitespace bool
	HasOpenAll            bool
	HasDefaultSearchAlias bool
	RegionOverride        *awsworkflow.Region
	RegionQuery           *string
	RemainingQuery        string
}

func (q *Query) IsEmpty() bool {
	return q.Service == nil && q.SubService == nil && q.RemainingQuery == "" && !q.HasOpenAll && q.RegionOverride == nil && q.RegionQuery == nil
}

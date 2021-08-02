package parsers

import (
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
)

type Query struct {
	RawQuery              string
	Service               *awsworkflow.AwsService
	SubService            *awsworkflow.AwsService
	HasTrailingWhitespace bool
	HasOpenAll            bool
	HasDefaultSearchAlias bool
	regionOverride        *awsconfig.Region
	RegionQuery           *string
	ProfileOverride       *awsconfig.Profile
	ProfileQuery          *string
	RemainingQuery        string
}

func (q *Query) IsEmpty() bool {
	return strings.Trim(q.RawQuery, " ") == ""
}

func (q *Query) GetRegionOverride() *awsconfig.Region {
	if q.Service == nil || !q.Service.HasGlobalRegion {
		return q.regionOverride
	}
	return nil
}

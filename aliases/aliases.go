package aliases

import "os"

var OverrideAwsRegion string
var OverrideAwsProfile string
var Search string

func init() {
	OverrideAwsRegion = os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_OVERRIDE_AWS_REGION_ALIAS")
	if OverrideAwsRegion == "" {
		OverrideAwsRegion = "$"
	}

	OverrideAwsProfile = os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_OVERRIDE_AWS_PROFILE_ALIAS")
	if OverrideAwsProfile == "" {
		OverrideAwsProfile = "@"
	}

	Search = os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_SEARCH_ALIAS")
	if Search == "" {
		Search = ","
	}
}

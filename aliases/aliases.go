package aliases

import "os"

var OverrideAwsRegion string
var Search string

func init() {
	OverrideAwsRegion = os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_OVERRIDE_AWS_REGION_ALIAS")
	if OverrideAwsRegion == "" {
		OverrideAwsRegion = "$"
	}

	Search = os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_SEARCH_ALIAS")
	if Search == "" {
		Search = ","
	}
}

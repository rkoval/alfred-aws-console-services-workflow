package parsers

import (
	"net/url"
	"strings"
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
	"github.com/stretchr/testify/assert"
)

func TestParseConsoleServicesYml(t *testing.T) {
	awsworkflow.InitAWSConsoleDomain("us-west-2")
	awsServices := ParseConsoleServicesYml("../console-services.yml")
	for _, awsService := range awsServices {
		if !awsService.HasGlobalRegion && !strings.HasPrefix(awsService.Url, "https://") {
			assert.Containsf(t, awsService.Url, "#", "url must contain '#' for region query parameter expansion to work properly")
		}

		region := ""
		if !awsService.HasGlobalRegion {
			region = "us-west-2"
		}
		rawUrl := util.ConstructAWSConsoleUrl(awsService.Url, region)
		_, err := url.Parse(rawUrl)
		assert.NoError(t, err)
	}
}

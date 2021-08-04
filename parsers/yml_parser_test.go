package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConsoleServicesYml(t *testing.T) {
	awsServices := ParseConsoleServicesYml("../console-services.yml")
	for _, awsService := range awsServices {
		if !awsService.HasGlobalRegion && len(awsService.SubServices) > 0 {
			assert.Containsf(t, awsService.Url, "#", "url must contain '#' for region query parameter expansion to work properly")
		}
	}
}

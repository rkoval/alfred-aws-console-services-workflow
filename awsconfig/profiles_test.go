package awsconfig

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
)

func TestGetAwsProfiles(t *testing.T) {
	cwd, _ := os.Getwd()
	credentialsPath := filepath.Join(cwd, "../tests/test_aws_credentials_file")
	configPath := filepath.Join(cwd, "../tests/test_aws_config_file")

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", credentialsPath)
	os.Setenv("AWS_CONFIG_FILE", configPath)

	awsProfiles = nil
	profiles := GetAwsProfiles()

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	createAndSaveSnapshot(t, profiles)
	verifyProfiles(t, profiles, credentialsPath, configPath)
}

func createAndSaveSnapshot(t *testing.T, profiles []Profile) {
	snap := make([]map[string]string, len(profiles))
	for i, p := range profiles {
		snap[i] = map[string]string{
			"name":   p.Name,
			"region": p.Region,
		}
	}
	cupaloy.SnapshotT(t, snap)
}

func verifyProfiles(t *testing.T, profiles []Profile, credentialsPath, configPath string) {
	credentialsIniFile, err := ini.Load(credentialsPath)
	assert.NoError(t, err)

	configIniFile, err := ini.Load(configPath)
	assert.NoError(t, err)

	profileNames := getProfileNames(profiles)
	expectedProfiles := getExpectedProfiles(credentialsIniFile, configIniFile)

	for source, profileMap := range expectedProfiles {
		for name := range profileMap {
			assert.Contains(t, profileNames, name,
				"Expected profile %s from %s to be loaded", name, source)
		}
	}

	verifyProfileRegions(t, profiles, configIniFile)

	// Implementation-specific behavior: profiles only in config file aren't loaded
	assert.NotContains(t, profileNames, "profile2")
}

func getExpectedProfiles(credentialsFile, configFile *ini.File) map[string]map[string]bool {
	result := map[string]map[string]bool{
		"credentials": make(map[string]bool),
		"sso":         make(map[string]bool),
	}

	// Get profiles from credentials file
	for _, section := range credentialsFile.Sections() {
		if section.Name() != ini.DefaultSection {
			result["credentials"][section.Name()] = true
		}
	}

	// Get SSO profiles from config file
	for _, section := range configFile.Sections() {
		if section.Name() != ini.DefaultSection && strings.HasPrefix(section.Name(), "profile ") {
			profileName := strings.TrimPrefix(section.Name(), "profile ")

			ssoSessionKey, _ := section.GetKey("sso_session")
			ssoStartUrlKey, _ := section.GetKey("sso_start_url")

			if ssoSessionKey != nil || ssoStartUrlKey != nil {
				result["sso"][profileName] = true
			}
		}
	}

	return result
}

func verifyProfileRegions(t *testing.T, profiles []Profile, configFile *ini.File) {
	for _, profile := range profiles {
		sectionName := profile.Name
		if sectionName != "default" {
			sectionName = "profile " + sectionName
		}

		section, err := configFile.GetSection(sectionName)
		if err != nil {
			continue
		}

		regionKey, err := section.GetKey("region")
		if err != nil {
			continue
		}

		assert.Equal(t, regionKey.Value(), profile.Region,
			"Profile %s should have region %s from config",
			profile.Name, regionKey.Value())
	}
}

func getProfileNames(profiles []Profile) []string {
	names := make([]string, len(profiles))
	for i, profile := range profiles {
		names[i] = profile.Name
	}
	return names
}

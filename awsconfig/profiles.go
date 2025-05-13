package awsconfig

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

// Profile represents an AWS profile with its configuration
type Profile struct {
	Name   string
	Region string
}

// Package-level variables
var awsProfiles []Profile

// GetAwsProfiles returns all AWS profiles from credentials and config files
func GetAwsProfiles() []Profile {
	if len(awsProfiles) <= 0 {
		loadAwsProfiles()
	}
	return awsProfiles
}

// GetAwsCredentialsFilePath returns the path to AWS credentials file
func GetAwsCredentialsFilePath() string {
	path := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if path == "" {
		path = config.DefaultSharedCredentialsFilename()
	}
	return path
}

// GetAwsProfileFilePath returns the path to AWS config file
func GetAwsProfileFilePath() string {
	path := os.Getenv("AWS_CONFIG_FILE")
	if path == "" {
		path = config.DefaultSharedConfigFilename()
	}
	return path
}

// loadAwsProfiles loads profiles from AWS credentials and config files
func loadAwsProfiles() {
	credentialsIniFile, err := ini.Load(GetAwsCredentialsFilePath())
	if err != nil {
		log.Println(err)
	}
	configIniFile, err := ini.Load(GetAwsProfileFilePath())
	if err != nil {
		log.Println(err)
	}

	awsProfiles = nil

	// First load profiles from credentials file
	if credentialsIniFile != nil {
		for _, section := range credentialsIniFile.Sections() {
			if section.Name() == ini.DefaultSection {
				continue
			}

			profileName := section.Name()
			var region string

			if configIniFile != nil {
				var configSectionName string
				if profileName == "default" {
					configSectionName = profileName
				} else {
					configSectionName = "profile " + profileName
				}
				configSection, _ := configIniFile.GetSection(configSectionName)
				region = getRegionFromSection(configSection)
			}

			addProfile(profileName, region)
		}
	}

	// Then load SSO profiles from config file
	if configIniFile != nil {
		for _, section := range configIniFile.Sections() {
			if section.Name() == ini.DefaultSection || !strings.HasPrefix(section.Name(), "profile ") {
				continue
			}

			profileName := strings.TrimPrefix(section.Name(), "profile ")
			if profileExists(profileName) {
				continue
			}

			ssoSessionKey, _ := section.GetKey("sso_session")
			ssoStartUrlKey, _ := section.GetKey("sso_start_url")

			if ssoSessionKey != nil || ssoStartUrlKey != nil {
				region := getRegionFromSection(section)
				addProfile(profileName, region)
			}
		}
	}
}

// addProfile creates and adds a profile to awsProfiles
func addProfile(name, region string) {
	profile := Profile{
		Name:   name,
		Region: region,
	}
	awsProfiles = append(awsProfiles, profile)
}

// profileExists checks if a profile with the given name already exists in awsProfiles
func profileExists(name string) bool {
	for _, p := range awsProfiles {
		if p.Name == name {
			return true
		}
	}
	return false
}

// getRegionFromSection extracts the region value from a section if it exists
func getRegionFromSection(section *ini.Section) string {
	if section == nil {
		return ""
	}

	regionKey, _ := section.GetKey("region")
	if regionKey != nil {
		return regionKey.Value()
	}
	return ""
}

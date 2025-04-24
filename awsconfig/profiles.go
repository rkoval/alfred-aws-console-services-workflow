package awsconfig

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

type Profile struct {
	Name   string
	Region string
}

var awsProfiles []Profile

func GetAwsProfiles() []Profile {
	if len(awsProfiles) <= 0 {
		loadAwsProfiles()
	}
	return awsProfiles
}

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
	if credentialsIniFile != nil {
		for _, section := range credentialsIniFile.Sections() {
			if section.Name() == ini.DefaultSection {
				// AWS does not specify anything useful in the root section
				continue
			}
			profile := Profile{
				Name: section.Name(),
			}

			if configIniFile != nil {
				var sectionName string
				if profile.Name == "default" {
					sectionName = profile.Name
				} else {
					// all other non-"default" profiles have special prefix
					sectionName = "profile " + profile.Name
				}
				configProfileSection, _ := configIniFile.GetSection(sectionName)
				if configProfileSection != nil {
					regionKey, _ := configProfileSection.GetKey("region")
					if regionKey != nil {
						profile.Region = regionKey.Value()
					}
				}
			}
			awsProfiles = append(awsProfiles, profile)
		}
	}
}

func GetAwsCredentialsFilePath() string {
	// see https://docs.aws.amazon.com/sdkref/latest/guide/file-location.html
	path := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if path == "" {
		path = config.DefaultSharedCredentialsFilename()
	}
	return path
}

func GetAwsProfileFilePath() string {
	// see https://docs.aws.amazon.com/sdkref/latest/guide/file-location.html
	path := os.Getenv("AWS_CONFIG_FILE")
	if path == "" {
		path = config.DefaultSharedConfigFilename()
	}
	return path
}

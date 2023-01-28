package parsers

import (
	"log"
	"os"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"gopkg.in/yaml.v2"
)

func ParseConsoleServicesYml(ymlPath string) []awsworkflow.AwsService {
	awsServices := []awsworkflow.AwsService{}
	yamlFile, err := os.ReadFile(ymlPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &awsServices)
	if err != nil {
		log.Fatal(err)
	}
	return awsServices
}

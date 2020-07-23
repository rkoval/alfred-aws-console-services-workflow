package parsers

import (
	"io/ioutil"
	"log"

	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"gopkg.in/yaml.v2"
)

func ParseConsoleServicesYml() []core.AwsService {
	awsServices := []core.AwsService{}
	yamlFile, err := ioutil.ReadFile("../console-services.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &awsServices)
	if err != nil {
		log.Fatal(err)
	}
	return awsServices
}

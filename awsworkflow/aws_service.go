package awsworkflow

import "github.com/aws/aws-sdk-go-v2/aws"

type AwsService struct {
	Id               string       `yaml:"id"`
	Name             string       `yaml:"name"`
	ShortName        string       `yaml:"short_name"`
	Description      string       `yaml:"description"`
	Url              string       `yaml:"url"`
	HomeID           string       `yaml:"home_id"`
	ExtraSearchTerms []string     `yaml:"extra_search_terms"`
	SubServices      []AwsService `yaml:"sub_services"`
	HasGlobalRegion  bool         `yaml:"has_global_region"`
}

func (this AwsService) GetName() string {
	if this.ShortName != "" {
		return this.ShortName + " – " + this.Name
	}
	return this.Name
}

func (this AwsService) GetRegion(cfg aws.Config) string {
	region := cfg.Region
	if this.HasGlobalRegion {
		region = ""
	}
	return region
}

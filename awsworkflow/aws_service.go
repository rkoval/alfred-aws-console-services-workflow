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

func (s *AwsService) GetName() string {
	if s.ShortName != "" {
		return s.ShortName + " – " + s.Name
	}
	return s.Name
}

func (s *AwsService) GetRegion(cfg aws.Config) string {
	if s.HasGlobalRegion {
		return ""
	}
	return cfg.Region
}

func (s *AwsService) HasSubServices() bool {
	return len(s.SubServices) > 0
}

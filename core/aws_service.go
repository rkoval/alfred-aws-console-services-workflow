package core

type AwsService struct {
	Id               string              `yaml:"id"`
	Name             string              `yaml:"name"`
	ShortName        string              `yaml:"short_name"`
	Description      string              `yaml:"description"`
	Url              string              `yaml:"url"`
	ExtraSearchTerms []string            `yaml:"extra_search_terms"`
	Sections         []AwsServiceSection `yaml:"sections"`
}

func (this AwsService) GetName() string {
	if this.ShortName != "" {
		return this.ShortName
	}
	return this.Name
}

type AwsServiceSection struct {
	Id          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Url         string `yaml:"url"`
}

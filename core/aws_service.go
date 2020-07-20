package core

type AwsService struct {
	Id               string   `yaml:"id"`
	Name             string   `yaml:"name"`
	ShortName        string   `yaml:"short_name"`
	Description      string   `yaml:"description"`
	Url              string   `yaml:"url"`
	ExtraSearchTerms []string `yaml:"extra_search_terms"`
}

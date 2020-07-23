package core

type AwsService struct {
	Id               string       `yaml:"id"`
	Name             string       `yaml:"name"`
	ShortName        string       `yaml:"short_name"`
	Description      string       `yaml:"description"`
	Url              string       `yaml:"url"`
	ExtraSearchTerms []string     `yaml:"extra_search_terms"`
	SubServices      []AwsService `yaml:"sub_services"`
}

func (this AwsService) GetName() string {
	if this.ShortName != "" {
		return this.ShortName + " – " + this.Name
	}
	return this.Name
}

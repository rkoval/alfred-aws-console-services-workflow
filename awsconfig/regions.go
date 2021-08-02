package awsconfig

type Region struct {
	Name        string
	Description string
}

var AllAWSRegions []Region = []Region{
	{
		Name:        "us-east-2",
		Description: "US East (Ohio)",
	},
	{
		Name:        "us-east-1",
		Description: "US East (N. Virginia)",
	},
	{
		Name:        "us-west-1",
		Description: "US West (N. California)",
	},
	{
		Name:        "us-west-2",
		Description: "US West (Oregon)",
	},
	{
		Name:        "af-south-1",
		Description: "Africa (Cape Town)",
	},
	{
		Name:        "ap-east-1",
		Description: "Asia Pacific (Hong Kong)",
	},
	{
		Name:        "ap-south-1",
		Description: "Asia Pacific (Mumbai)",
	},
	{
		Name:        "ap-northeast-3",
		Description: "Asia Pacific (Osaka)",
	},
	{
		Name:        "ap-northeast-2",
		Description: "Asia Pacific (Seoul)",
	},
	{
		Name:        "ap-southeast-1",
		Description: "Asia Pacific (Singapore)",
	},
	{
		Name:        "ap-southeast-2",
		Description: "Asia Pacific (Sydney)",
	},
	{
		Name:        "ap-northeast-1",
		Description: "Asia Pacific (Tokyo)",
	},
	{
		Name:        "ca-central-1",
		Description: "Canada (Central)",
	},
	{
		Name:        "cn-north-1",
		Description: "China (Beijing)",
	},
	{
		Name:        "cn-northwest-1",
		Description: "China (Ningxia)",
	},
	{
		Name:        "eu-central-1",
		Description: "Europe (Frankfurt)",
	},
	{
		Name:        "eu-west-1",
		Description: "Europe (Ireland)",
	},
	{
		Name:        "eu-west-2",
		Description: "Europe (London)",
	},
	{
		Name:        "eu-south-1",
		Description: "Europe (Milan)",
	},
	{
		Name:        "eu-west-3",
		Description: "Europe (Paris)",
	},
	{
		Name:        "eu-north-1",
		Description: "Europe (Stockholm)",
	},
	{
		Name:        "me-south-1",
		Description: "Middle East (Bahrain)",
	},
	{
		Name:        "sa-east-1",
		Description: "South America (São Paulo)",
	},
}

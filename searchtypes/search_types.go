package searchtypes

type SearchType int

const (
	None SearchType = iota
	Services
	SubServices

	EC2Instances
	EC2SecurityGroups
	ElasticBeanstalkEnvironments
	S3Buckets
)

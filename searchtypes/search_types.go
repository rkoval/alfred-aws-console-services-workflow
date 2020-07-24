package searchtypes

type SearchType int

const (
	Services SearchType = iota + 1
	SubServices

	EC2Instances
	EC2SecurityGroups
	ElasticBeanstalkEnvironments
	S3Buckets
)

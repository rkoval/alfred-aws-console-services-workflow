module github.com/rkoval/alfred-aws-console-services-workflow

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.7.1
	github.com/aws/aws-sdk-go-v2/config v1.4.1
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.6.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.5.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.11.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.8.0
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.5.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.5.0
	github.com/aws/aws-sdk-go-v2/service/lambda v1.4.0
	github.com/aws/aws-sdk-go-v2/service/rds v1.5.0
	github.com/aws/aws-sdk-go-v2/service/route53 v1.7.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.11.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.6.0
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.6.0
	github.com/aws/smithy-go v1.6.0
	github.com/bradleyjkemp/cupaloy v2.3.0+incompatible
	github.com/cheekybits/genny v1.0.0
	github.com/deanishe/awgo v0.25.0
	github.com/dnaeon/go-vcr v1.0.1
	github.com/iancoleman/strcase v0.1.3
	github.com/stretchr/testify v1.6.1
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

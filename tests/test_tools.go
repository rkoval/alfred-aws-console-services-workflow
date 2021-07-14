package tests

import (
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
)

func NewAWSRecorderSession(fixtureName string) (aws.Config, *recorder.Recorder) {
	var mode recorder.Mode
	if os.Getenv("RECORD_VCR") != "" {
		mode = recorder.ModeRecording
	} else {
		mode = recorder.ModeReplaying
	}
	r, err := recorder.NewAsMode(fixtureName+"_aws_fixture", mode, nil)
	if err != nil {
		panic(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Request.Headers, "X-Amz-Date")
		delete(i.Request.Headers, "X-Amz-Content-Sha256")
		delete(i.Request.Headers, "User-Agent")
		delete(i.Request.Headers, "Amz-Sdk-Request")
		delete(i.Request.Headers, "Amz-Sdk-Invocation-Id")
		delete(i.Response.Headers, "X-Amzn-Requestid")
		delete(i.Response.Headers, "Date")
		delete(i.Response.Headers, "X-Amz-Id-2")
		delete(i.Response.Headers, "X-Amz-Request-Id")
		delete(i.Response.Headers, "Content-Length")
		return nil
	})

	r.AddFilter(func(i *cassette.Interaction) error {
		i.Response.Body = sanitizeBody(i.Response.Body)
		return nil
	})

	cfg := awsworkflow.InitAWS(r)

	return cfg, r
}

func PanicOnError(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}

var environmentIdRegex *regexp.Regexp = regexp.MustCompile(`e-[a-zA-Z0-9]{8,}`)
var instanceIdRegex *regexp.Regexp = regexp.MustCompile(`i-[a-zA-Z0-9]{8,}`)
var dbIdRegex *regexp.Regexp = regexp.MustCompile(`db-[a-zA-Z0-9]{8,}`)
var amiIdRegex *regexp.Regexp = regexp.MustCompile(`ami-[a-zA-Z0-9]{8,}`)
var vpcIdRegex *regexp.Regexp = regexp.MustCompile(`vpc-[a-zA-Z0-9]{8,}`)
var subnetIdRegex *regexp.Regexp = regexp.MustCompile(`subnet-[a-zA-Z0-9]{8,}`)
var securityGroupIdRegex *regexp.Regexp = regexp.MustCompile(`sg-[a-zA-Z0-9]{8,}`)
var expandedSecurityGroupIdRegex *regexp.Regexp = regexp.MustCompile(`securitygroup-[a-zA-Z0-9]{8,}`)
var volumeIdRegex *regexp.Regexp = regexp.MustCompile(`vol-[a-zA-Z0-9]{8,}`)
var attachmentIdRegex *regexp.Regexp = regexp.MustCompile(`eni-attach-[a-zA-Z0-9]{8,}`)
var reservationIdRegex *regexp.Regexp = regexp.MustCompile(`r-[a-zA-Z0-9]{8,}`)

var accountIdInArn *regexp.Regexp = regexp.MustCompile(`:[0-9]{10,}:`)
var longNumberInXmlTag *regexp.Regexp = regexp.MustCompile(`>[0-9]{8,}<`) // we're going to assume that any numeric xml values are identifications of some sort, so just sanitize it
var uuidv2Regex *regexp.Regexp = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
var iso8601Regex *regexp.Regexp = regexp.MustCompile(`\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d+)?(([+-]\d\d:\d\d)|Z)?`)

var idTagRegex *regexp.Regexp = regexp.MustCompile(`<(id|ID|DbiResourceId|HostedZoneId)>.+</(id|ID|DbiResourceId|HostedZoneId)>`)
var keyNameTagRegex *regexp.Regexp = regexp.MustCompile(`<keyName>.+</keyName>`)
var masterUsernameTagRegex *regexp.Regexp = regexp.MustCompile(`<MasterUsername>.+</MasterUsername>`)

var ipv4Regex *regexp.Regexp = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
var macAddressRegex *regexp.Regexp = regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)

var beanstalkSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBSecurityGroup-[0-9A-Z]{10,}`)
var beanstalkLoadBalancerSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBLoadBalancerSecurityGroup-[0-9A-Z]{10,}`)
var beanstalkAutoScalingGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBAutoScalingGroup-[0-9A-Z]{10,}`)

var amazonawsUrlRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9]+\.amazonaws\.com`)
var beanstalkUrlSubdomainRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.elasticbeanstalk\.com`)
var internalUrlRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9]+\.internal`)

func sanitizeBody(body string) string {
	body = uuidv2Regex.ReplaceAllString(body, "00000000-0000-0000-0000-000000000000")
	body = environmentIdRegex.ReplaceAllString(body, "e-aaaaaaaaaa")
	body = instanceIdRegex.ReplaceAllString(body, "i-aaaaaaaaaa")
	body = dbIdRegex.ReplaceAllString(body, "db-AAAAAAAAAA")
	body = amiIdRegex.ReplaceAllString(body, "ami-aaaaaaaaaa")
	body = vpcIdRegex.ReplaceAllString(body, "vpc-aaaaaaaaaa")
	body = subnetIdRegex.ReplaceAllString(body, "subnet-aaaaaaaaaa")
	body = securityGroupIdRegex.ReplaceAllString(body, "sg-aaaaaaaaaa")
	body = expandedSecurityGroupIdRegex.ReplaceAllString(body, "securitygroup-aaaaaaaaaa")
	body = volumeIdRegex.ReplaceAllString(body, "vol-aaaaaaaaaa")
	body = attachmentIdRegex.ReplaceAllString(body, "eni-attach-aaaaaaaaaa")
	body = reservationIdRegex.ReplaceAllString(body, "r-aaaaaaaaaa")

	body = accountIdInArn.ReplaceAllString(body, ":0000000000:")
	body = longNumberInXmlTag.ReplaceAllString(body, ">00000000<")
	body = iso8601Regex.ReplaceAllString(body, "2020-01-01T00:00:00.000Z")

	body = idTagRegex.ReplaceAllString(body, "<id>000000000000</id>")
	body = masterUsernameTagRegex.ReplaceAllString(body, "<MasterUsername>aaaaaaaaaaaa</MasterUsername>")
	body = keyNameTagRegex.ReplaceAllString(body, "<keyName>aaaaaaaaaa</keyName>")

	body = ipv4Regex.ReplaceAllString(body, "0.0.0.0")
	body = macAddressRegex.ReplaceAllString(body, "00:00:00:00:00:00")

	body = beanstalkSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBSecurityGroup-AAAAAAAAAAAA")
	body = beanstalkLoadBalancerSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBLoadBalancerSecurityGroup-AAAAAAAAAAAA")
	body = beanstalkAutoScalingGroupNameRegex.ReplaceAllString(body, "AWSEBAutoScalingGroup-AAAAAAAAAAAA")

	body = amazonawsUrlRegex.ReplaceAllString(body, "subdomain.us-west-2.service.amazonaws.com")
	body = beanstalkUrlSubdomainRegex.ReplaceAllString(body, "subdomain.us-west-2.elasticbeanstalk.com")
	body = internalUrlRegex.ReplaceAllString(body, "subdomain.us-west-2.service.internal")

	return body
}

package tests

import (
	"os"
	"regexp"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func NewAWSRecorder(fixtureLocation string) *recorder.Recorder {
	var mode recorder.Mode
	if os.Getenv("RECORD_VCR") != "" {
		mode = recorder.ModeRecording
	} else {
		mode = recorder.ModeReplaying
	}
	r, err := recorder.NewAsMode(fixtureLocation, mode, nil)
	if err != nil {
		panic(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Request.Headers, "X-Amz-Date")
		delete(i.Response.Headers, "X-Amzn-Requestid")
		delete(i.Response.Headers, "Date")
		return nil
	})

	r.AddFilter(func(i *cassette.Interaction) error {
		i.Response.Body = sanitizeBody(i.Response.Body)
		return nil
	})

	return r
}

var environmentIdRegex *regexp.Regexp = regexp.MustCompile(`e-[a-zA-Z0-9]{8,}`)
var vpcIdRegex *regexp.Regexp = regexp.MustCompile(`vpc-[a-zA-Z0-9]{8,}`)
var securityGroupIdRegex *regexp.Regexp = regexp.MustCompile(`sg-[a-zA-Z0-9]{8,}`)
var accountIdInArn *regexp.Regexp = regexp.MustCompile(`:[0-9]{10,}:`)
var longNumberInXmlTag *regexp.Regexp = regexp.MustCompile(`>[0-9]{8,}<`) // we're going to assume that any numeric xml values are identifications of some sort, so just sanitize it
var uuidv2Regex *regexp.Regexp = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
var ipv4Regex *regexp.Regexp = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
var beanstalkSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBSecurityGroup-[0-9A-Z]{10,}`)
var beanstalkLoadBalancerSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBLoadBalancerSecurityGroup-[0-9A-Z]{10,}`)
var elbUrlSubdomain *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.elb.amazonaws\.com`)
var beanstalkUrlSubdomain *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.elasticbeanstalk\.com`)

func sanitizeBody(body string) string {
	body = environmentIdRegex.ReplaceAllString(body, "e-aaaaaaaaaa")
	body = vpcIdRegex.ReplaceAllString(body, "vpc-aaaaaaaaaa")
	body = securityGroupIdRegex.ReplaceAllString(body, "sg-aaaaaaaaaa")
	body = accountIdInArn.ReplaceAllString(body, ":0000000000:")
	body = longNumberInXmlTag.ReplaceAllString(body, ">00000000<")
	body = uuidv2Regex.ReplaceAllString(body, "00000000-0000-0000-0000-000000000000")
	body = ipv4Regex.ReplaceAllString(body, "0.0.0.0")
	body = beanstalkSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBSecurityGroup-AAAAAAAAAAAA")
	body = beanstalkLoadBalancerSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBLoadBalancerSecurityGroup-AAAAAAAAAAAA")
	body = elbUrlSubdomain.ReplaceAllString(body, "subdomain.us-west-2.elb.amazonaws.com")
	body = beanstalkUrlSubdomain.ReplaceAllString(body, "subdomain.us-west-2.elasticbeanstalk.com")
	return body
}

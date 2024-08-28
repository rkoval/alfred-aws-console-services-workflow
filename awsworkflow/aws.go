package awsworkflow

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

var SigninToken string

func InitAWS(transport http.RoundTripper, profile *awsconfig.Profile, region *awsconfig.Region, wf *aw.Workflow) aws.Config {
	cfg, err := awsconfig.GetAwsConfig(profile, region)
	if err != nil {
		panic(err)
	}

	if transport != nil {
		cfg.HTTPClient = &http.Client{
			Transport: transport,
		}
	}

	InitAWSConsoleDomain(cfg.Region)
	err = InitAWSSession(cfg, wf)
	if err != nil {
		panic(err)
	}

	return cfg
}

var defaultAwsConsoleDomain string = "console.aws.amazon.com"
var defaultAwsConsoleDomainChina string = "console.amazonaws.cn"
var defaultAwsConsoleDomainUsGov string = "console.amazonaws-us-gov.com"

func InitAWSConsoleDomain(region string) {
	oldAwsConsoleDomain := os.Getenv("ALRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_CONSOLE_DOMAIN")
	if oldAwsConsoleDomain != "" {
		panic(errors.New("`ALRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_CONSOLE_DOMAIN` env var was renamed to `ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_CONSOLE_DOMAIN` due to the misspelling. Please update your config"))
	}
	awsConsoleDomain := os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_CONSOLE_DOMAIN")
	if awsConsoleDomain == "" {
		if strings.HasPrefix(region, "cn-") {
			awsConsoleDomain = defaultAwsConsoleDomainChina
		} else if strings.HasPrefix(region, "us-gov-") {
			awsConsoleDomain = defaultAwsConsoleDomainUsGov
		} else {
			awsConsoleDomain = defaultAwsConsoleDomain
		}
	}
	util.AWSConsoleDomain = awsConsoleDomain
}

func InitAWSSession(cfg aws.Config, wf *aw.Workflow) error {
	fileName := "aws_session_" + util.GetProfile(cfg) + ".json"
	var awsSession util.SessionJSON

	// Attempt to load session from cache
	if err := loadSessionFromCache(wf, fileName, &awsSession); err != nil {
		// Cache miss or error, create a new session
		var err error
		awsSession, err = util.CreateSession(cfg)
		if err != nil {
			return err
		}
		if err := wf.Cache.StoreJSON(fileName, awsSession); err != nil {
			return err
		}
	}

	// Check if session is expired
	if isSessionExpired(awsSession) {
		awsSession, err := util.CreateSession(cfg)
		if err != nil {
			return err
		}
		if err := wf.Cache.StoreJSON(fileName, awsSession); err != nil {
			return err
		}
	}

	util.SigninToken = awsSession.SigninToken
	return nil

}

func loadSessionFromCache(wf *aw.Workflow, fileName string, awsSession *util.SessionJSON) error {
	if wf.Cache.Exists(fileName) {
		if err := wf.Cache.LoadJSON(fileName, awsSession); err != nil {
			return err
		}
		return nil
	}
	return errors.New("cache does not exist")
}

func isSessionExpired(session util.SessionJSON) bool {
	return session.Expires.Before(time.Now())
}

func GetImageIcon(id string) *aw.Icon {
	icon := &aw.Icon{Value: "images/" + id + ".png"}
	return icon
}

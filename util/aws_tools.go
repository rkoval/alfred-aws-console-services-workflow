package util

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cloudformationTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudwatchlogsTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elasticacheTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elasticbeanstalkTypes "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// this will be set by init
var AWSConsoleDomain string
var SigninToken string

func ConstructAWSConsoleUrl(path, region string) string {
	if strings.HasPrefix(path, "http") {
		// some URLs are just global landing pages, so just return it as-is if we detect a protocol
		return path
	}

	if AWSConsoleDomain == "" {
		panic(errors.New("AWSConsoleDomain was not initialized"))
	}

	// TODO append region query param dynamically here to avoid page redirections and facilitate faster loading

	var urlBuilder strings.Builder
	urlBuilder.WriteString("https://")
	if region != "" {
		urlBuilder.WriteString(region)
		urlBuilder.WriteString(".")
	}

	urlBuilder.WriteString(AWSConsoleDomain)

	if region != "" {
		urlBuilder.WriteString(strings.Replace(path, "#", "?region="+region+"#", 1))
	} else {
		urlBuilder.WriteString(path)
	}

	url := GenerateLogoutLoginUrl(urlBuilder.String())

	return url.String()
}

func GetProfile(cfg aws.Config) string {
	for _, configSource := range cfg.ConfigSources {
		switch cs := configSource.(type) {
		case config.SharedConfig:
			return cs.Profile
		}
	}
	return ""
}

func GetEC2TagValue(tags []ec2Types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func GetEndOfArn(arn string) string {
	splitArn := strings.Split(arn, ":")
	return splitArn[len(splitArn)-1]
}

func GetCloudFormationTagValue(tags []cloudformationTypes.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func GetEC2InstanceStateEmoji(instanceState ec2Types.InstanceState) string {
	switch name := instanceState.Name; name {
	case ec2Types.InstanceStateNameRunning:
		return "🟢"
	case ec2Types.InstanceStateNameShuttingDown:
		return "🟡"
	case ec2Types.InstanceStateNameStopping:
		return "🟡"
	case ec2Types.InstanceStateNameStopped:
		return "🔴"
	case ec2Types.InstanceStateNameTerminated:
		return "🔴"
	case ec2Types.InstanceStateNamePending:
		return "⚪️"
	}

	return "❔"
}

func GetElasticBeanstalkHealthEmoji(environmentHealth elasticbeanstalkTypes.EnvironmentHealth) string {
	switch environmentHealth {
	case elasticbeanstalkTypes.EnvironmentHealthGreen:
		return "🟢"
	case elasticbeanstalkTypes.EnvironmentHealthYellow:
		return "🟡"
	case elasticbeanstalkTypes.EnvironmentHealthRed:
		return "🔴"
	case elasticbeanstalkTypes.EnvironmentHealthGrey:
		return "⚪️"
	}

	return "❔"
}

func GetElasticacheCacheClusterSubtitle(entity elasticacheTypes.CacheCluster) string {
	engineArray := []string{}
	engineArray = AppendString(engineArray, entity.Engine)
	engineArray = AppendString(engineArray, entity.EngineVersion)
	engineString := strings.Join(engineArray, " ")
	subtitleArray := []string{}
	subtitleArray = AppendString(subtitleArray, &engineString)
	subtitleArray = AppendString(subtitleArray, entity.CacheNodeType)
	subtitleArray = AppendString(subtitleArray, entity.CacheClusterStatus)
	subtitle := strings.Join(subtitleArray, " – ")

	return subtitle
}

func awsOuterEscape(s string) string {
	// QueryEscape doesn't escape tildes like browsers do, so we must do that here manually
	// see https://github.com/golang/go/issues/47379
	return url.QueryEscape(strings.ReplaceAll(url.QueryEscape(s), "~", "%7E"))
}

func awsInnerEscape(s string) string {
	// QueryEscape doesn't escape tildes like browsers do, so we must do that here manually
	// see https://github.com/golang/go/issues/47379
	return strings.ReplaceAll(url.QueryEscape(s), "%", "*")
}

func ConstructCloudwatchInsightsQueryDetail(entity cloudwatchlogsTypes.QueryDefinition) string {
	// cloudwatch insights has a crazy url scheme for referencing queries instead of just by ID, so we must reconstruct that here.
	// logic below adapted from https://stackoverflow.com/questions/60796991/is-there-a-way-to-generate-the-aws-console-urls-for-cloudwatch-log-group-filters
	// TODO do outer escaping manually for perf?
	var queryDetailBuilder strings.Builder
	queryDetailBuilder.WriteString(awsOuterEscape("~(end~0~start~-3600~timeType~'RELATIVE~unit~'seconds~editorString~'"))
	queryDetailBuilder.WriteString(awsInnerEscape(*entity.QueryString))
	queryDetailBuilder.WriteString(awsOuterEscape("~isLiveTail~false~queryId~'"))
	queryDetailBuilder.WriteString(awsInnerEscape(*entity.QueryDefinitionId))
	queryDetailBuilder.WriteString(awsOuterEscape("~source~("))
	for _, logGroupName := range entity.LogGroupNames {
		queryDetailBuilder.WriteString(awsOuterEscape("~'"))
		queryDetailBuilder.WriteString(awsInnerEscape(logGroupName))
	}
	queryDetailBuilder.WriteString(awsOuterEscape("))"))
	return strings.ReplaceAll(queryDetailBuilder.String(), "%", "$")
}

func GenerateLogoutLoginUrl(destinationUrl string) url.URL {
	loginUrlParams := url.Values{}
	loginUrlParams.Add("Action", "login")
	loginUrlParams.Add("Issuer", "aws-vault")
	loginUrlParams.Add("SigninToken", SigninToken)
	loginUrlParams.Add("Destination", destinationUrl)

	loginUrl := url.URL{
		Scheme:   "https",
		Host:     "us-east-1.signin.aws.amazon.com",
		Path:     "/federation",
		RawQuery: url.Values(loginUrlParams).Encode(),
	}

	logoutUrlParams := url.Values{}
	logoutUrlParams.Add("Action", "logout")
	logoutUrlParams.Add("redirect_uri", loginUrl.String())

	finalUrl := url.URL{
		Scheme:   "https",
		Host:     "signin.aws.amazon.com",
		Path:     "/oauth",
		RawQuery: url.Values(logoutUrlParams).Encode(),
	}

	if err := validateAWSUrl(finalUrl.String()); err != nil {
		panic(err)
	}

	return finalUrl
}

func validateAWSUrl(url string) error {
	// TODO
	return nil
}

type SessionJSON struct {
	Expires         time.Time
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	SigninToken     string
}

func CreateSession(cfg aws.Config) (SessionJSON, error) {
	credentials, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		return SessionJSON{}, err
	}

	session := map[string]string{
		"sessionId":    credentials.AccessKeyID,
		"sessionKey":   credentials.SecretAccessKey,
		"sessionToken": credentials.SessionToken,
	}

	if session["sessionToken"] == "" {
		session, err = fetchFederatedSession(cfg)
		if err != nil {
			return SessionJSON{}, err
		}
	}

	signinToken, err := fetchSigninToken(session)
	if err != nil {
		return SessionJSON{}, err
	}

	return SessionJSON{
		// Login links expire after 15 minutes, so we'll use that as our session expiration
		Expires:         time.Now().Add(time.Second * 900),
		AccessKeyId:     session["sessionId"],
		SecretAccessKey: session["sessionKey"],
		SessionToken:    session["sessionToken"],
		SigninToken:     signinToken,
	}, nil
}

func fetchFederatedSession(cfg aws.Config) (map[string]string, error) {
	client := sts.NewFromConfig(cfg)

	input := &sts.GetFederationTokenInput{
		DurationSeconds: aws.Int32(43200),
		Name:            aws.String("alfred-aws"),
		Policy:          aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`),
	}

	output, err := client.GetFederationToken(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"sessionId":    *output.Credentials.AccessKeyId,
		"sessionKey":   *output.Credentials.SecretAccessKey,
		"sessionToken": *output.Credentials.SessionToken,
	}, nil
}

func fetchSigninToken(session map[string]string) (string, error) {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	signinTokenURL := url.URL{
		Scheme: "https",
		Host:   "signin.aws.amazon.com",
		Path:   "/federation",
	}

	params := url.Values{
		"Action":          {"getSigninToken"},
		"Session":         {string(sessionJSON)},
		"SessionDuration": {"43200"},
	}

	signinTokenURL.RawQuery = params.Encode()

	resp, err := http.Get(signinTokenURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var signinTokenJSON map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&signinTokenJSON); err != nil {
		return "", err
	}

	signinToken := signinTokenJSON["SigninToken"]
	if signinToken == "" {
		return "", errors.New("SigninToken not found in response")
	}

	return signinToken, nil
}

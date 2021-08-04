package util

import (
	"errors"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cloudformationTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudwatchlogsTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elasticacheTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elasticbeanstalkTypes "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
)

// this will be set by init
var AWSConsoleDomain string

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
	return urlBuilder.String()
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

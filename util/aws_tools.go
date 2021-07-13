package util

import (
	"strings"

	cloudformationTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elasticbeanstalkTypes "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
)

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

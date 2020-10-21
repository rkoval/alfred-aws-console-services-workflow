package util

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

func GetEC2TagValue(tags []*ec2.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func GetCloudFormationTagValue(tags []*cloudformation.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func GetEC2InstanceStateEmoji(instanceState ec2.InstanceState) string {
	switch name := *instanceState.Name; name {
	case ec2.InstanceStateNameRunning:
		return "🟢"
	case ec2.InstanceStateNameShuttingDown:
		return "🟡"
	case ec2.InstanceStateNameStopping:
		return "🟡"
	case ec2.InstanceStateNameStopped:
		return "🔴"
	case ec2.InstanceStateNameTerminated:
		return "🔴"
	case ec2.InstanceStateNamePending:
		return "⚪️"
	}

	return "❔"
}

func GetElasticBeanstalkHealthEmoji(environmentHealth string) string {
	switch environmentHealth {
	case elasticbeanstalk.EnvironmentHealthGreen:
		return "🟢"
	case elasticbeanstalk.EnvironmentHealthYellow:
		return "🟡"
	case elasticbeanstalk.EnvironmentHealthRed:
		return "🔴"
	case elasticbeanstalk.EnvironmentHealthGrey:
		return "⚪️"
	}

	return "❔"
}

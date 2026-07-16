package awsconfig

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig/providers"
)

func GetAwsConfig(profile *Profile, region *Region) (aws.Config, error) {
	providerType := os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_AUTH_PROVIDER")

	profileLoadOptionsFunc, regionLoadOptionsFunc, err := getLoadOptions(profile, region)
	if err != nil {
		return aws.Config{}, err
	}

	switch providerType {
	case "aws-vault":
		return getAWSVaultConfig(profile, profileLoadOptionsFunc, regionLoadOptionsFunc)
	default:
		return getDefaultConfig(profileLoadOptionsFunc, regionLoadOptionsFunc)
	}
}

func getLoadOptions(profile *Profile, region *Region) (func(*config.LoadOptions) error, func(*config.LoadOptions) error, error) {
	profileLoadOptionsFunc := func(o *config.LoadOptions) error { return nil }
	regionLoadOptionsFunc := func(o *config.LoadOptions) error { return nil }

	if profile != nil {
		profileLoadOptionsFunc = config.WithSharedConfigProfile(profile.Name)
		if profile.Region != "" {
			regionLoadOptionsFunc = config.WithRegion(profile.Region)
		}
	}
	if region != nil {
		regionLoadOptionsFunc = config.WithRegion(region.Name)
	}

	return profileLoadOptionsFunc, regionLoadOptionsFunc, nil
}

func getAWSVaultConfig(profile *Profile, profileLoadOptionsFunc, regionLoadOptionsFunc func(*config.LoadOptions) error) (aws.Config, error) {
	if profile == nil || profile.Name == "" {
		profile = &Profile{Name: "default"}
	}

	log.Println("using auth provider aws-vault with profile", profile.Name)

	provider := providers.NewAWSVaultCredentialsProvider(profile.Name)
	cfg, err := config.LoadDefaultConfig(context.TODO(), profileLoadOptionsFunc, regionLoadOptionsFunc, config.WithCredentialsProvider(provider.Cache))
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS SDK config: %v", err)
	}

	return cfg, nil
}

func getDefaultConfig(profileLoadOptionsFunc, regionLoadOptionsFunc func(*config.LoadOptions) error) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		profileLoadOptionsFunc,
		regionLoadOptionsFunc,
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load default AWS SDK config: %v", err)
	}

	return cfg, nil
}

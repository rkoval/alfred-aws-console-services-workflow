package providers

import (
	"bufio"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type AWSVaultCredentialsProvider struct {
	Profile string
	Cache   *aws.CredentialsCache
}

func NewAWSVaultCredentialsProvider(profile string) *AWSVaultCredentialsProvider {
	provider := &AWSVaultCredentialsProvider{
		Profile: profile,
	}

	provider.Cache = aws.NewCredentialsCache(provider)

	return provider
}

func (c AWSVaultCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	output, err := executeAwsVaultCommand(c.Profile)
	if err != nil {
		return aws.Credentials{}, errors.New("failed to execute aws-vault command: " + err.Error())
	}

	keyValueMap, err := parseKeyValuePairs(output)
	if err != nil {
		return aws.Credentials{}, errors.New("failed to parse aws-vault output: " + err.Error())
	}

	creds, err := mapToCredentials(keyValueMap)
	if err != nil {
		return aws.Credentials{}, errors.New("failed to map credentials: " + err.Error())
	}

	return creds, nil
}

func executeAwsVaultCommand(profile string) (string, error) {
	cmd := exec.Command("aws-vault", "exec", profile, "--", "env")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func parseKeyValuePairs(output string) (map[string]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	m := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return m, errors.New("failed to parse line: " + line)
		}
		m[parts[0]] = parts[1]
	}

	return m, nil
}

const (
	accessKeyIDKey       = "AWS_ACCESS_KEY_ID"
	secretAccessKeyKey   = "AWS_SECRET_ACCESS_KEY"
	sessionTokenKey      = "AWS_SESSION_TOKEN"
	sessionExpirationKey = "AWS_CREDENTIAL_EXPIRATION"
	timeFormat           = "2006-01-02T15:04:05Z"
)

func mapToCredentials(m map[string]string) (aws.Credentials, error) {
	var creds aws.Credentials
	var err error

	creds.AccessKeyID = m[accessKeyIDKey]
	creds.SecretAccessKey = m[secretAccessKeyKey]
	creds.SessionToken = m[sessionTokenKey]
	creds.CanExpire = true
	creds.Expires, err = time.Parse(timeFormat, m[sessionExpirationKey])

	if err != nil {
		return creds, errors.New("failed to parse expiration time")
	}

	return creds, nil
}

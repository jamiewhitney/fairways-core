package secretstore

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AWSSecretStore struct {
	sm *secretsmanager.SecretsManager
}

func NewAWSSecretManager(ctx context.Context) (SecretStore, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &AWSSecretStore{
		secretsmanager.New(sess),
	}, nil
}

func (as AWSSecretStore) GetSecret(ctx context.Context, key string) (string, error) {
	var secretID, versionID, versionStage string

	current := &secretID
	for _, ch := range key {
		if ch == '@' {
			current = &versionID
			continue
		}

		if ch == '#' {
			current = &versionStage
			continue
		}

		*current += string(ch)
	}

	req := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	if versionID != "" {
		req.VersionId = aws.String(versionID)
	}

	if versionStage != "" {
		req.VersionStage = aws.String(versionStage)
	}

	result, err := as.sm.GetSecretValueWithContext(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret %v: %w", key, err)
	}

	if v := aws.StringValue(result.SecretString); v != "" {
		return v, nil
	}

	return string(result.SecretBinary), nil
}

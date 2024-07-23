package secretstore

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-envconfig"
	"strings"
)

type Config struct {
	Type string `env:"SECRET_MANAGER" json:"type"`
}

type SecretStore interface {
	GetSecret(ctx context.Context, key string) (string, error)
}

func GetSecretStore(ctx context.Context, config *Config) (SecretStore, error) {
	switch config.Type {
	case "vault":
		return NewHashiCorpVault(ctx)
	case "aws":
		return NewAWSSecretManager(ctx)
	default:
		return nil, fmt.Errorf("unknown secret manager type: %s", config.Type)
	}
}

func Resolver(sm SecretStore) envconfig.MutatorFunc {
	return func(ctx context.Context, originalKey string, resolvedKey string, originalValue string, currentValue string) (newValue string, stop bool, err error) {
		s, err := resolve(sm, originalValue)
		if err != nil {
			return "", false, err
		}
		return s, false, nil
	}
}

func resolve(sm SecretStore, secretRef string) (string, error) {
	secretPrefix := "secret://"
	if !strings.HasPrefix(secretRef, secretPrefix) {
		return secretRef, nil
	}

	secretRef = strings.TrimPrefix(secretRef, secretPrefix)

	return sm.GetSecret(context.Background(), secretRef)
}

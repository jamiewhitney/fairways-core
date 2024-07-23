package config

//import (
//	"context"
//	"errors"
//	"github.com/fairways/fairways-core/pkg/cache"
//	"github.com/fairways/fairways-core/pkg/database"
//	"github.com/fairways/fairways-core/pkg/secretstore"
//	"github.com/hashicorp/vault/api"
//	"github.com/sethvargo/go-envconfig"
//	"os"
//	"time"
//)
//
//type AppConfig struct {
//	Environment string `env:"env"`
//	Platform    string `env:"platform"`
//
//	Database database.Config
//	Cache    cache.Config
//
//	DBVaultSecret string `env:"DB_VAULT_SECRET"`
//
//	HTTPPort         string `env:"HTTP_PORT, default=3000"`
//	HTTPReadTimeout  time.Duration
//	HTTPWriteTimeout time.Duration
//	HTTPIdleTimeout  time.Duration
//	GRPCPort         string `env:"GRPC_PORT,  default=3001"`
//
//	VaultAddr     string `env:"VAULT_ADDR,default=http://localhost:8200"`
//	VaultAuthRole string `env:"VAULT_AUTH_ROLE"`
//}
//
//func LoadConfig() (*AppConfig, error) {
//	var c AppConfig
//
//	env, ok := os.LookupEnv("ENV")
//	if !ok {
//		return nil, errors.New("ENV not set")
//	}
//	switch env {
//	case "development":
//		if err := envconfig.Process(context.TODO(), &c); err != nil {
//			return nil, err
//		}
//	case "production":
//		if err := envconfig.Process(context.TODO(), &c); err != nil {
//			return nil, err
//		}
//
//		secret, err := secretstore.NewHashiCorpVault(&secretstore.HashiCorpVaultConfig{
//			Role:        c.VaultAuthRole,
//		})
//		if err != nil {
//			return nil, err
//		}
//		if secret == nil {
//			return nil, api.ErrSecretNotFound
//		}
//
//
//		return nil, err
//	default:
//		return nil, errors.New("ENV not set")
//	}
//
//	return &c, nil
//}

package pricing

import (
	"github.com/jamiewhitney/fairways-core/pkg/cache"
	"github.com/jamiewhitney/fairways-core/pkg/database"
	"github.com/jamiewhitney/fairways-core/pkg/secretstore"
)

type Config struct {
	Database database.Config                   `json:"database"`
	Secrets  *secretstore.HashiCorpVaultConfig `json:"secrets"`
	Cache    cache.Config                      `json:"cache"`

	Port string `env:"PORT, default=3000"`
}

func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}

func (c *Config) CacheConfig() *cache.Config {
	return &c.Cache
}

func (c *Config) VaultConfig() *secretstore.HashiCorpVaultConfig {
	return c.Secrets
}

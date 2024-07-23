package catalog

import (
	"github.com/jamiewhitney/fairways-core/pkg/cache"
	"github.com/jamiewhitney/fairways-core/pkg/database"
	internal_services "github.com/jamiewhitney/fairways-core/pkg/internal-services"
	"github.com/jamiewhitney/fairways-core/pkg/pubsub"
	"github.com/jamiewhitney/fairways-core/pkg/secretstore"
)

type Config struct {
	Database database.Config    `json:"database"`
	Secrets  secretstore.Config `json:"secrets"`
	Cache    cache.Config       `json:"cache"`
	Pubsub   pubsub.Config      `json:"pubsub"`

	teeTimeServiceConfig internal_services.ServiceConfig
	pricingServiceConfig internal_services.ServiceConfig
	bookingServiceConfig internal_services.ServiceConfig

	Port string `env:"PORT,default=3000"`
}

func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}

func (c *Config) CacheConfig() *cache.Config {
	return &c.Cache
}

func (c *Config) PubsubConfig() *pubsub.Config {
	return &c.Pubsub
}

func (c *Config) SecretStoreConfig() *secretstore.Config {
	return &c.Secrets
}

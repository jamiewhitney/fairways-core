package tee_time

import (
	"github.com/jamiewhitney/fairways-core/pkg/cache"
	"github.com/jamiewhitney/fairways-core/pkg/database"
	internal_services "github.com/jamiewhitney/fairways-core/pkg/internal-services"
	"github.com/jamiewhitney/fairways-core/pkg/pubsub"
	"github.com/jamiewhitney/fairways-core/pkg/secretstore"
	booking_pb "github.com/jamiewhitney/fairways-core/protobufs/booking"
	pricing_pb "github.com/jamiewhitney/fairways-core/protobufs/pricing"
)

type Config struct {
	Database database.Config                  `json:"database"`
	Secrets  secretstore.HashiCorpVaultConfig `json:"secrets"`
	Cache    cache.Config                     `json:"cache"`
	Pubsub   pubsub.Config                    `json:"pubsub"`

	PricingServiceClient pricing_pb.PricingServiceClient
	PricingConf          internal_services.ServiceConfig
	BookingServiceClient booking_pb.BookingServiceClient
	BookingConf          internal_services.ServiceConfig

	Port string `json:"port" env:"PORT, default=3000"`
}

func (c *Config) VaultConfig() *secretstore.HashiCorpVaultConfig {
	return &c.Secrets
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

func (c *Config) PricingServiceConfig() internal_services.ServiceConfig {
	return c.PricingConf
}

func (c *Config) BookingServiceConfig() internal_services.ServiceConfig {
	return c.BookingConf
}

func (c *Config) PricingClient() pricing_pb.PricingServiceClient {
	return c.PricingServiceClient
}

func (c *Config) BookingClient() booking_pb.BookingServiceClient {
	return c.BookingServiceClient
}

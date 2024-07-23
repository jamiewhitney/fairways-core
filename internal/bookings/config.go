package bookings

import (
	"github.com/jamiewhitney/fairways-core/pkg/database"
)

type Config struct {
	Database database.Config `json:"database"`
	Port     string          `env:"PORT, default=3000"`
}

func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}

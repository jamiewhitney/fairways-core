package database

import (
	"time"
)

type Config struct {
	Name               string        `env:"DB_NAME" json:",omitempty"`
	UserPass           string        `env:"DB_USERPASS,required" json:",omitempty"`
	Host               string        `env:"DB_HOST, default=localhost" json:"host,omitempty"`
	Port               string        `env:"DB_PORT, default=3306" json:"port,omitempty"`
	SSLMode            string        `env:"DB_SSLMODE, default=require" json:"ssl_mode,omitempty"`
	ConnectionTimeout  int           `env:"DB_CONNECT_TIMEOUT" json:",omitempty"`
	SSLCertPath        string        `env:"DB_SSLCERT" json:",omitempty"`
	SSLKeyPath         string        `env:"DB_SSLKEY" json:",omitempty"`
	SSLRootCertPath    string        `env:"DB_SSLROOTCERT" json:",omitempty"`
	PoolMinConnections string        `env:"DB_POOL_MIN_CONNS" json:",omitempty"`
	PoolMaxConnections string        `env:"DB_POOL_MAX_CONNS" json:",omitempty"`
	PoolMaxConnLife    time.Duration `env:"DB_POOL_MAX_CONN_LIFETIME, default=5m" json:",omitempty"`
	PoolMaxConnIdle    time.Duration `env:"DB_POOL_MAX_CONN_IDLE_TIME, default=1m" json:",omitempty"`
}

func (c *Config) DatabaseConfig() *Config {
	return c
}
